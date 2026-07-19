package ticket

import (
	"go-learning/internal/dto"
	"go-learning/internal/service"
	"go-learning/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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

// Order Ticket By User
// @Summary      Order Ticket By User
// @Description  Order Ticket By User
// @Tags         vetaute management
// @Accept       json
// @Produce      json
// @Param        payload  body  dto.OrderRequest true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /tickets/order [post]
func (c *cTicketItem) PlaceOrderByUser(ctx *gin.Context) {
	validation, exists := ctx.Get("validation")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Validator not found"})
	}

	var body dto.OrderRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}

	// check validation
	err := validation.(*validator.Validate).Struct(body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response.SuccessResponse(ctx, response.ErrCodeSuccess, body)
}
