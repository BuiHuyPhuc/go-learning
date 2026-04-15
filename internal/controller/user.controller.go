package controller

import (
	"go-learning/internal/service"
	res "go-learning/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: service.NewUserService(),
	}
}

func (uc *UserController) GetUserById(c *gin.Context) {
	res.SuccessResponse(c, res.ErrCodeSuccess, uc.userService.GetInfoUser())
	// if err != nil {
	//   return res.ErrorResponse(c, res.ErrCodeParamInvalid, "email is invalid")
	// }
}
