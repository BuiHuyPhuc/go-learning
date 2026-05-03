package repo

import (
	"go-learning/global"
	"go-learning/internal/database"
)

// type UserRepo struct {
// }

// func NewUserRepo() *UserRepo {
// 	return &UserRepo{}
// }

// func (ur *UserRepo) GetInfoUser() string {
//   return "Phúc nè"
// }

// INTERFACE
type IUserRepository interface {
	GetUserByEmail(email string) bool
}

type userRepository struct {
	sqlc *database.Queries
}

// GetUserByEmail implements [IUserRepository].
func (ur *userRepository) GetUserByEmail(email string) bool {
	user, err := ur.sqlc.GetUserByEmailSQLC(ctx, email)
	if err != nil {
		return false
	}
	return user.UsrID != 0
}

func NewUserRepository() IUserRepository {
	return &userRepository{
		sqlc: database.New(global.Mdbc),
	}
}
