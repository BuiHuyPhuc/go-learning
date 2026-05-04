package account

import (
	"go-learning/internal/service"
	"go-learning/pkg/response"

	"github.com/gin-gonic/gin"
)

// management controller Login User
var Login = new(cUserLogin)

type cUserLogin struct{}

func (c *cUserLogin) Login(ctx *gin.Context) {
	err := service.UserLogin().Login(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}

	response.SuccessResponse(ctx, response.ErrCodeSuccess, nil)
}
