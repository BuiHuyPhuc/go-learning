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
	"sync"
)

type sTicketItem struct {
	r                *database.Queries
	distributedCache service.IDistributedCache
	localCache       service.IRistrettoCache
}

func NewTicketItemImpl(r *database.Queries, distributedCache service.IDistributedCache) *sTicketItem {
	localCache, err := NewRistrettoCacheImpl()
	if err != nil {
		fmt.Printf("Initialize local cache failed: %v\n", err)
	}

	return &sTicketItem{r, distributedCache, localCache}
}

var mu sync.Mutex

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

	// 1. get ticket item from local cache
	out, err = s.getTicketItemFromLocalCache(ctx, ticketId, "v1")
	if err != nil {
		return out, fmt.Errorf("error getting ticket item from local cache: %w", err)
	}
	if (out != dto.TicketItemResponse{}) {
		fmt.Println("CALL service GetTicketItemById from LOCALCACHE with ", ticketId)
		return out, nil
	}

	// lấy khóa để tránh lượng đồng thời cao cùng 1 thời điểm gây áp lực xuống db
	// mu.Lock()
	// defer mu.Unlock()

	// 2. get ticket item from distributed cache
	out, err = s.getTicketItemFromDistributedCache(ctx, ticketId)
	if err != nil {
		return out, fmt.Errorf("error getting ticket item from distributed cache: %w", err)
	}
	if (out != dto.TicketItemResponse{}) {
		fmt.Println("CALL service GetTicketItemById from REDIS with ", ticketId)
		return out, nil
	}

	// 3. get ticket item from database
	out, err = s.getTicketItemFromDatabaseLock(ctx, ticketId)
	if err != nil {
		return out, err
	}

	fmt.Println("CALL service GetTicketItemById from MYSQL with ", ticketId)
	return out, nil
}

func (s *sTicketItem) getTicketItemFromDatabaseLock(ctx context.Context, ticketId int) (out dto.TicketItemResponse, err error) {
	lockKey := fmt.Sprintf("lock:ticketItem:%d", ticketId)
	err = s.distributedCache.WithDistributedLock(ctx, lockKey, 5, func(ctx context.Context) error {
		fmt.Printf("LOCK  ACQUIRED -> QUERY DATABASE BY TICKETID %d\n", ticketId)

		out, err = s.getTicketItemFromDatabase(ctx, ticketId)
		if err != nil {
			return err
		}

		return nil
	})

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

	// put to distributed cache
	err = s.distributedCache.Set(ctx, s.getKeyTicketItemCache(ticketId), ticketItemCacheJSON, consts.TIME_OTP_REGISTER*60)
	if err != nil {
		return out, fmt.Errorf("save redis failed: %v", err)
	}

	// put to local cache
	isSuccess := s.localCache.SetWithTTL(ctx, s.getKeyTicketItemCache(ticketId), ticketItem)
	if !isSuccess {
		return out, fmt.Errorf("save localcache failed: %v", err)
	}

	return mapper.ToTicketItemDTO(ticketItem), nil
}

func (s *sTicketItem) getTicketItemFromDistributedCache(ctx context.Context, ticketId int) (out dto.TicketItemResponse, err error) {
	ticketItemCache, err := s.distributedCache.Get(ctx, s.getKeyTicketItemCache(ticketId))
	if err != nil {
		return out, fmt.Errorf("failed to get ticket item from distributed cache: %v", err)
	}

	if ticketItemCache == "" {
		return out, nil
	}

	if err := json.Unmarshal([]byte(ticketItemCache), &out); err != nil {
		return out, fmt.Errorf("parse redis data failed: %v", err)
	}

	// put to local cache
	s.localCache.SetWithTTL(ctx, s.getKeyTicketItemCache(ticketId), out)

	return out, nil
}

func (s *sTicketItem) getTicketItemFromLocalCache(ctx context.Context, ticketId int, version string) (out dto.TicketItemResponse, err error) {
	ticketItemCache, ok := s.localCache.Get(ctx, s.getKeyTicketItemCache(ticketId))
	if !ok {
		return out, fmt.Errorf("failed to get ticket item from local cache: %v", err)
	}

	ticketItemCacheStr, ok := ticketItemCache.(string)
	if !ok {
		return out, fmt.Errorf("ticketItemCache is not string")
	}

	if ticketItemCacheStr == "" {
		return out, nil
	}

	if err := json.Unmarshal([]byte(ticketItemCacheStr), &out); err != nil {
		return out, fmt.Errorf("parse redis data failed: %v", err)
	}

	return out, nil
}

func (s *sTicketItem) getKeyTicketItemCache(ticketId int) string {
	return fmt.Sprintf("ticketItem-%d", ticketId)
}
