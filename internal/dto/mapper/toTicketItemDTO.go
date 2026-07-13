package mapper

import (
	"go-learning/internal/database"
	"go-learning/internal/dto"
)

func ToTicketItemDTO(ticketItem database.GetTicketItemByIdRow) dto.TicketItemResponse {
	return dto.TicketItemResponse{
		TicketId:       int(ticketItem.ID),
		TicketName:     ticketItem.Name,
		SotckInitial:   int(ticketItem.StockInitial),
		StockAvailable: int(ticketItem.StockAvailable),
	}
}
