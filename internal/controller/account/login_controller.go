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

// User Login documentation
// @Summary      User Login
// @Description  User Login
// @Tags         account management
// @Accept       json
// @Produce      json
// @Param        payload  body  dto.LoginRequest true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /users/login [post]
func (c *cUserLogin) Login(ctx *gin.Context) {
	var params dto.LoginRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}

	codeStatus, result, err := service.UserLogin().Login(ctx, &params)
	if err != nil {
		response.ErrorResponse(ctx, codeStatus, err.Error())
		return
	}

	response.SuccessResponse(ctx, codeStatus, result)
}

// User Registration documentation
// @Summary      User Registration
// @Description  When user is registerd send otp  to email
// @Tags         account management
// @Accept       json
// @Produce      json
// @Param        payload  body  dto.RegisterRequest true "payload"
// @Success      200  {object}  response.ResponseData
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

// Verify OTP Login By User
// @Summary      Verify OTP Login By User
// @Description  When user is registerd send otp  to email
// @Tags         account management
// @Accept       json
// @Produce      json
// @Param        payload  body  dto.VerifyOTPRequest true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /users/verify-account [post]
func (c *cUserLogin) VerifyOTP(ctx *gin.Context) {
	var params dto.VerifyOTPRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}

	result, err := service.UserLogin().VerifyOTP(ctx, &params)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrInvalidOTP, err.Error())
		return
	}

	response.SuccessResponse(ctx, response.ErrCodeSuccess, result)
}

// Update Password Register
// @Summary      Update Password Register
// @Description  When user is registerd send otp  to email
// @Tags         account management
// @Accept       json
// @Produce      json
// @Param        payload  body  dto.UpdatePasswordRegisterRequest true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /users/update-pass-register [post]
func (c *cUserLogin) UpdatePasswordRegister(ctx *gin.Context) {
	var params dto.UpdatePasswordRegisterRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}

	result, err := service.UserLogin().UpdatePasswordRegister(ctx, params.Token, params.Password)
	if err != nil {
		response.ErrorResponse(ctx, result, err.Error())
		return
	}

	response.SuccessResponse(ctx, response.ErrCodeSuccess, result)
}
