package repo

import (
	"fmt"
	"go-learning/global"
	"time"
)

type IUserAuthRepository interface {
	AddOTP(email string, otp int, expiration int64) error
}

type UserAuthRepository struct {
}

func (uar *UserAuthRepository) AddOTP(email string, otp int, expiration int64) error {
	key := fmt.Sprintf("usr:%s:opt", email) // usr:email:otp
	return global.Rdb.Set(ctx, key, otp, time.Duration(expiration)).Err()
}

func NewUserAuthRepository() IUserAuthRepository {
	return &UserAuthRepository{}
}
