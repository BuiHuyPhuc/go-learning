package account

import (
	"go-learning/global"
	"go-learning/internal/dto"
	"go-learning/internal/service"
	"go-learning/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

// User Registration documentation
// @Summary      User Registration
// @Description  When user is registerd send otp  to email
// @Tags         account management
// @Accept       json
// @Produce      json
// @Param        payload  body  dto.RegisterRequest true "payload"
// @Success      200  {object}  response.Response
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /users/register [post]
func (c *cUserLogin) Register(ctx *gin.Context) {
	var params dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}

	codeStatus, err := service.UserLogin().Register(ctx, &params)
	if err != nil {
		global.Logger.Error("Error registering user OTP", zap.Error(err))
		response.ErrorResponse(ctx, codeStatus, err.Error())
		return
	}

	response.SuccessResponse(ctx, codeStatus, nil)
}
