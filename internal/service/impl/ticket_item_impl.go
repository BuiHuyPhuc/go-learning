package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"go-learning/internal/consts"
	"go-learning/internal/database"
	"go-learning/internal/dto"
	"go-learning/internal/dto/mapper"
	"go-learning/internal/service"
)

type sTicketItem struct {
	r                *database.Queries
	distributedCache service.IDistributedCache
}

func NewTicketItemImpl(r *database.Queries, distributedCache service.IDistributedCache) *sTicketItem {
	return &sTicketItem{r, distributedCache}
}

func (s *sTicketItem) GetTicketItemById(ctx context.Context, ticketId int) (out dto.TicketItemResponse, err error) {
	fmt.Println("CAL service GetTicketItemById...")
	// ticketItem, err := s.r.GetTicketItemById(ctx, int64(ticketId))
	// if err != nil {
	// 	return nil, err
	// }

	// out = &dto.TicketItemResponse{
	// 	TicketId:       int(ticketItem.ID),
	// 	TicketName:     ticketItem.Name,
	// 	StockAvailable: int(ticketItem.StockAvailable),
	// 	SotckInitial:   int(ticketItem.StockInitial),
	// }

	out, err = s.getTicketItemFromDistributedCache(ctx, ticketId)
	if err != nil {
		return out, fmt.Errorf("error getting ticket item from distributed cache: %w", err)
	}

	if (out != dto.TicketItemResponse{}) {
		fmt.Println("CALL service GetTicketItemById from REDIS with ", ticketId)
		return out, nil
	}

	out, err = s.getTicketItemFromDatabase(ctx, ticketId)
	if err != nil {
		return out, err
	}

	fmt.Println("CALL service GetTicketItemById from MYSQL with ", ticketId)
	return out, nil
}

func (s *sTicketItem) getTicketItemFromDatabase(ctx context.Context, ticketId int) (out dto.TicketItemResponse, err error) {
	ticketItem, err := s.r.GetTicketItemById(ctx, int64(ticketId))
	if err != nil {
		return out, err
	}

	ticketItemCacheJSON, err := json.Marshal(ticketItem)
	if err != nil {
		return out, fmt.Errorf("convert to json failed: %v", err)
	}

	err = s.distributedCache.Set(ctx, s.getKeyTicketItemCache(ticketId), ticketItemCacheJSON, consts.TIME_OTP_REGISTER*60)
	if err != nil {
		return out, fmt.Errorf("save redis failed: %v", err)
	}

	return mapper.ToTicketItemDTO(ticketItem), nil
}

func (s *sTicketItem) getTicketItemFromDistributedCache(ctx context.Context, ticketId int) (out dto.TicketItemResponse, err error) {
	ticketItemmCache, err := s.distributedCache.Get(ctx, s.getKeyTicketItemCache(ticketId))
	if err != nil {
		return out, fmt.Errorf("failed to get ticket item cache: %v", err)
	}

	if ticketItemmCache == "" {
		return out, nil
	}

	if err := json.Unmarshal([]byte(ticketItemmCache), &out); err != nil {
		return out, fmt.Errorf("parse redis data failed: %v", err)
	}

	return out, nil
}

func (s *sTicketItem) getKeyTicketItemCache(ticketId int) string {
	return fmt.Sprintf("ticketItem-%d", ticketId)
}
