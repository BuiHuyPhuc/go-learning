package ticket

import (
	"go-learning/internal/service"
	"go-learning/pkg/response"

	"github.com/gin-gonic/gin"
)

var TicketItem = new(cTicketItem)

type cTicketItem struct{}

func (c *cTicketItem) GetTicketItemById(ctx *gin.Context) {
	ticketItem, err := service.TicketItem().GetTicketItemById(ctx, 1)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}

	response.SuccessResponse(ctx, response.ErrCodeSuccess, ticketItem)
}
