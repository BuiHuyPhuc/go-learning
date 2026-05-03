package service

import (
	"fmt"
	"go-learning/internal/repo"
	"go-learning/internal/utils/crypto"
	"go-learning/internal/utils/random"
	"go-learning/internal/utils/sendto"
	"go-learning/pkg/response"
	"strconv"
	"time"
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
	userRepo     repo.IUserRepository
	userAuthRepo repo.IUserAuthRepository
}

// Register implements [IUserService].
func (us *userService) Register(email, purpose string) int {
	// 0. hashEmail
	hashEmail := crypto.GetHash(email)
	fmt.Printf("hashEmail: %s\n", hashEmail)

	// 5. check OTP is available

	// 6. user spam ...

	// 1. check email exists in db
	if us.userRepo.GetUserByEmail(email) {
		return response.ErrCodeUserHasExists
	}

	// 2. new OTP > ...
	otp := random.GenerateSixDigitIOtp()
	if purpose == "test" {
		otp = 123456
	}
	fmt.Printf("Otp is: %d\n", otp)

	// 3. save OTP in Redis with expiration time
	err := us.userAuthRepo.AddOTP(hashEmail, otp, int64(10*time.Minute))
	if err != nil {
		return response.ErrInvalidOTP
	}

	// 4. send Email OTP
	// err = sendto.SendTextEmailOtp([]string{email}, "buihuyphuc6101997@gmail.com", strconv.Itoa(otp))
	err = sendto.SendTemplateEmailOtp(
		[]string{email},
		"buihuyphuc6101997@gmail.com",
		"otp-auth.html",
		map[string]interface{}{"otp": strconv.Itoa(otp)},
	)
	if err != nil {
		return response.ErrSendEmailOTP
	}

	// send OTP via Kafka
	// body := make(map[string]interface{})
	// body["otp"] = otp
	// body["email"] = email

	// jsonBody, _ := json.Marshal(body)
	// msg := kafka.Message{
	// 	Key:   []byte("otp-auth"),
	// 	Value: []byte(jsonBody),
	// 	Time:  time.Now(),
	// }
	// err = global.KafkaProducer.WriteMessages(context.Background(), msg)
	// if err != nil {
	// 	return response.ErrSendEmailOTP
	// }

	return response.ErrCodeSuccess
}

func NewUserService(
	userRepo repo.IUserRepository,
	userAuthRepo repo.IUserAuthRepository,
) IUserService {
	return &userService{
		userRepo,
		userAuthRepo,
	}
}
