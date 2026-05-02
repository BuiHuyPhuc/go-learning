//go:build wireinject

package wire

import (
	"go-learning/internal/controller"
	"go-learning/internal/repo"
	"go-learning/internal/service"

	"github.com/google/wire"
)

// cd internal/wire/ -> wire
func InitUserRouterHandler() (*controller.UserController, error) {
	wire.Build(
		repo.NewUserRepository,
		service.NewUserService,
		controller.NewUserController,
	)

	return new(controller.UserController), nil
}
