package service

import (
	"context"
	"go-learning/internal/dto"
)

type (
	//.. interface
	IUserAdmin interface {
		RemoveUser(ctx context.Context) error
		FindOneUser(ctx context.Context) error
	}

	IUserInfo interface {
		GetInfoByUserId(ctx context.Context) error
		GetAllUser(ctx context.Context) error
	}

	IUserLogin interface {
		Login(ctx context.Context, in *dto.LoginRequest) (statusCode int, out dto.LoginResponse, err error)
		Register(ctx context.Context, in *dto.RegisterRequest) (statusCode int, err error)
		VerifyOTP(ctx context.Context, in *dto.VerifyOTPRequest) (out dto.VerifyOTPResponse, err error)
		UpdatePasswordRegister(ctx context.Context, token string, password string) (userId int, err error)
	}
)

var (
	localUserAdmin IUserAdmin
	localUserInfo  IUserInfo
	localUserLogin IUserLogin
)

func UserAdmin() IUserAdmin {
	if localUserAdmin == nil {
		panic("Implement localUserAdmin not found for interface IUserAdmin")
	}
	return localUserAdmin
}
func InitUserAdmin(i IUserAdmin) { localUserAdmin = i }

func UserInfo() IUserInfo {
	if localUserInfo == nil {
		panic("Implement localUserInfo not found for interface IUserInfo")
	}
	return localUserInfo
}
func InitUserInfo(i IUserInfo) { localUserInfo = i }

func UserLogin() IUserLogin {
	if localUserLogin == nil {
		panic("Implement localUserLogin not found for interface IUserLogin")
	}
	return localUserLogin
}
func InitUserLogin(i IUserLogin) { localUserLogin = i }
