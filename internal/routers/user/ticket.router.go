package user

import (
	"go-learning/internal/controller/ticket"

	"github.com/gin-gonic/gin"
)

type TicketRouter struct{}

func (tr *TicketRouter) InitTicketRouter(Router *gin.RouterGroup) {
	// public router
	ticketRouterPublic := Router.Group("/tickets")
	{
		ticketRouterPublic.GET("/detail/:id", ticket.TicketItem.GetTicketItemById)
	}
}
