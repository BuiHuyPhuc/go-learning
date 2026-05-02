package controller

import (
	"go-learning/internal/service"
	"go-learning/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.IUserService
}

func NewUserController(
	userService service.IUserService,
) *UserController {
	return &UserController{
		userService,
	}
}

func (uc *UserController) Register(c *gin.Context) {
	result := uc.userService.Register(c.Query("email"), c.Query("purpose"))
	response.SuccessResponse(c, result, nil)
}

// func (uc *UserController) GetUserById(c *gin.Context) {
// 	res.SuccessResponse(c, res.ErrCodeSuccess, uc.userService.GetInfoUser())
// 	// if err != nil {
// 	//   return res.ErrorResponse(c, res.ErrCodeParamInvalid, "email is invalid")
// 	// }
// }
