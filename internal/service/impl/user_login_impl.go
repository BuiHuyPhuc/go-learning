package impl

import (
	"context"
	"database/sql"
	"fmt"
	"go-learning/global"
	"go-learning/internal/consts"
	"go-learning/internal/database"
	"go-learning/internal/dto"
	"go-learning/internal/utils"
	"go-learning/internal/utils/crypto"
	"go-learning/internal/utils/random"
	"go-learning/internal/utils/sendto"
	"go-learning/pkg/response"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type sUserLogin struct {
	r *database.Queries
}

func NewUserLoginImpl(r *database.Queries) *sUserLogin {
	return &sUserLogin{r}
}

func (s *sUserLogin) Login(ctx context.Context) error {
	return nil
}

func (s *sUserLogin) Register(ctx context.Context, in *dto.RegisterRequest) (statusCode int, err error) {
	// 1. hash email
	fmt.Printf("VerifyKey: %s\n", in.VerifyKey)
	fmt.Printf("VerifyType: %d\n", in.VerifyType)
	hashKey := crypto.GetHash(strings.ToLower(in.VerifyKey))
	fmt.Printf("hashKey: %s\n", hashKey)

	// 2. check user exists in user base
	userFound, err := s.r.CheckuserBaseExists(ctx, in.VerifyKey)
	if err != nil {
		return response.ErrCodeUserHasExists, err
	}

	if userFound > 0 {
		return response.ErrCodeUserHasExists, fmt.Errorf("user has already registered")
	}

	// 3. create OTP
	userKey := utils.GetUserKey(hashKey)
	otpFound, err := global.Rdb.Get(ctx, userKey).Result()

	switch {
	case err == redis.Nil:
		fmt.Println("Key does not exist")
	case err != nil:
		fmt.Println("Get failed:", err)
		return response.ErrInvalidOTP, err
	case otpFound != "":
		return response.ErrCodeOtpNotExists, fmt.Errorf("otp exists but not registered")
	}

	// 4. generate OTP
	otpNew := random.GenerateSixDigitIOtp()
	if in.VerifyPurpose == "TEST_USER" {
		otpNew = 123456
	}
	fmt.Printf("Otp is :::%d\n", otpNew)

	// 5. save OTP in Redis with expiration time
	err = global.Rdb.Set(ctx, userKey, strconv.Itoa(otpNew), time.Minute*time.Duration(consts.TIME_OTP_REGISTER)).Err()
	if err != nil {
		return response.ErrInvalidOTP, err
	}

	// 6. Send OTP
	switch in.VerifyType {
	case consts.EMAIL:
		err = sendto.SendTextEmailOtp([]string{in.VerifyKey}, consts.HOST_EMAIL, strconv.Itoa(otpNew))
		if err != nil {
			return response.ErrSendEmailOTP, err
		}

		// 7. save OTP to
		result, err := s.r.InsertOTPVerify(ctx, database.InsertOTPVerifyParams{
			VerifyOtp:     strconv.Itoa(otpNew),
			VerifyType:    sql.NullInt32{Int32: 1, Valid: true},
			VerifyKey:     in.VerifyKey,
			VerifyKeyHash: hashKey,
		})
		if err != nil {
			return response.ErrSendEmailOTP, err
		}

		// 8. getLastId
		lasIdVerifyUser, err := result.LastInsertId()
		if err != nil {
			return response.ErrSendEmailOTP, err
		}
		log.Println("lasIdVerifyUser", lasIdVerifyUser)

		return response.ErrCodeSuccess, nil

	case consts.MOBILE:
		return response.ErrCodeSuccess, nil
	}

	return response.ErrCodeSuccess, nil
}

func (s *sUserLogin) VerifyOTP(ctx context.Context) error {
	return nil
}

func (s *sUserLogin) UpdatePasswordRegister(ctx context.Context) error {
	return nil
}
