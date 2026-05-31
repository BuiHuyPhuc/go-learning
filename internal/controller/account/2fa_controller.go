package account

import (
	"go-learning/internal/dto"
	"go-learning/internal/service"
	"go-learning/internal/utils/context"
	"go-learning/pkg/response"
	"log"

	"github.com/gin-gonic/gin"
)

var TwoFA = new(cUser2FA)

type cUser2FA struct{}

// User Setup Two Factor Authentication
// @Summary      User Setup Two Factor Authentication
// @Description  User Setup Two Factor Authentication
// @Tags         account 2fa
// @Accept       json
// @Produce      json
// @param        Authorization header string true "Authorization token"
// @Param        payload  body  dto.SetupTwoFactorAuthRequest true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /users/two-factor/setup [post]
func (c *cUser2FA) SetupTwoFactorAuth(ctx *gin.Context) {
	var params dto.SetupTwoFactorAuthRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}

	// get userId from uuid (token)
	userId, err := context.GetUserIdFromUUID(ctx.Request.Context())
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeTwoFactorAuthSetupFailed, "UserId is not valid")
		return
	}
	log.Println("userId: ", userId)

	params.UserId = int(userId)
	codeStatus, err := service.UserLogin().SetupTwoFactorAuth(ctx, &params)
	if err != nil {
		response.ErrorResponse(ctx, codeStatus, err.Error())
		return
	}

	response.SuccessResponse(ctx, codeStatus, nil)
}

// User Verify Two Factor Authentication
// @Summary      User Verify Two Factor Authentication
// @Description  User Verify Two Factor Authentication
// @Tags         account 2fa
// @Accept       json
// @Produce      json
// @param        Authorization header string true "Authorization token"
// @Param        payload  body  dto.VerifyTwoFactorAuthRequest true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /users/two-factor/verify [post]
func (c *cUser2FA) VerifyTwoFactorAuth(ctx *gin.Context) {
	var params dto.VerifyTwoFactorAuthRequest
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}

	// get userId from uuid (token)
	userId, err := context.GetUserIdFromUUID(ctx.Request.Context())
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeTwoFactorAuthSetupFailed, "UserId is not valid")
		return
	}
	log.Println("userId: ", userId)

	params.UserId = int(userId)
	codeStatus, err := service.UserLogin().VerifyTwoFactorAuth(ctx, &params)
	if err != nil {
		response.ErrorResponse(ctx, codeStatus, err.Error())
		return
	}

	response.SuccessResponse(ctx, codeStatus, nil)
}
