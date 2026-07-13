package service

import (
	"context"
	"go-learning/internal/dto"
)

type (
	ITicketHome interface {
	}
	ITicketItem interface {
		GetTicketItemById(ctx context.Context, ticketId int) (out dto.TicketItemResponse, err error)
	}
)

var (
	localTicketHome ITicketHome
	localTicketItem ITicketItem
)

func TicketHome() ITicketHome {
	if localTicketHome == nil {
		panic("Implement localTicketHome not found for interface ITicketHome")
	}
	return localTicketHome
}
func InitTicketHome(i ITicketHome) { localTicketHome = i }

func TicketItem() ITicketItem {
	if localTicketItem == nil {
		panic("Implement localTicketItem not found for interface ITicketItem")
	}
	return localTicketItem
}
func InitTicketItem(i ITicketItem) { localTicketItem = i }
