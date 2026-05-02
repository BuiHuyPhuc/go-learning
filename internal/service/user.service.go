package service

import (
	"go-learning/internal/repo"
	"go-learning/pkg/response"
)

// type UserService struct {
//   userRepo *repo.UserRepo
// }

// func NewUserService() *UserService {
// 	return &UserService{
//     userRepo: repo.NewUserRepo(),
//   }
// }

// func (us *UserService) GetInfoUser() string {
//   return us.userRepo.GetInfoUser()
// }

// INTERFACE
type IUserService interface {
	Register(email, purpose string) int
}

type userService struct {
	userRepo repo.IUserRepository
}

// Register implements [IUserService].
func (us *userService) Register(email, purpose string) int {
	// 1. Check email exists
	if us.userRepo.GetUserByEmail(email) {
		return response.ErrCodeUserHasExists
	}
	return response.ErrCodeSuccess
}

func NewUserService(
	userRepo repo.IUserRepository,
) IUserService {
	return &userService{
		userRepo,
	}
}
