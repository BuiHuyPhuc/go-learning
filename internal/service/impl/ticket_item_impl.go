package impl

import (
	"context"
	"fmt"
	"go-learning/internal/database"
	"go-learning/internal/dto"
)

type sTicketItem struct {
	r *database.Queries
}

func NewTicketItemImpl(r *database.Queries) *sTicketItem {
	return &sTicketItem{r}
}

func (s *sTicketItem) GetTicketItemById(ctx context.Context, ticketId int) (out *dto.TicketItemResponse, err error) {
	fmt.Println("CAL service GetTicketItemById...")
	ticketItem, err := s.r.GetTicketItemById(ctx, int64(ticketId))
	if err != nil {
		return nil, err
	}

	out = &dto.TicketItemResponse{
		TicketId:       int(ticketItem.ID),
		TicketName:     ticketItem.Name,
		StockAvailable: int(ticketItem.StockAvailable),
		SotckInitial:   int(ticketItem.StockInitial),
	}

	return out, nil
}
