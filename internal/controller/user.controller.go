package controller

import (
	"fmt"
	"go-learning/internal/dto"
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
	var params dto.UserRegistratorRequest
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(c, response.ErrCodeParamInvalid, err.Error())
		return
	}

	fmt.Printf("Email params: %s\n", params.Email)
	result := uc.userService.Register(params.Email, params.Purpose)
	response.SuccessResponse(c, result, nil)
}

// func (uc *UserController) GetUserById(c *gin.Context) {
// 	res.SuccessResponse(c, res.ErrCodeSuccess, uc.userService.GetInfoUser())
// 	// if err != nil {
// 	//   return res.ErrorResponse(c, res.ErrCodeParamInvalid, "email is invalid")
// 	// }
// }
