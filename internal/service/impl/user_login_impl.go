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

func (s *sUserLogin) VerifyOTP(ctx context.Context, in *dto.VerifyOTPRequest) (out dto.VerifyOTPResponse, err error) {
	hashKey := crypto.GetHash(strings.ToLower(in.VerifyKey))

	userKey := utils.GetUserKey(hashKey)
	otpFound, err := global.Rdb.Get(ctx, userKey).Result()
	if err != nil {
		return out, err
	}

	if in.VerifyCode != otpFound {
		// Nếu như sai 3 lần trong vòng 1 phút?
		return out, fmt.Errorf("OTP not match")
	}

	// get OTP
	infoOTP, err := s.r.GetInfoOTP(ctx, hashKey)
	if err != nil {
		return out, err
	}

	// update status verified
	err = s.r.UpdateUserVerificationStatus(ctx, hashKey)
	if err != nil {
		return out, err
	}

	out.Token = infoOTP.VerifyKeyHash // token template
	out.Message = "success"

	return out, nil
}

func (s *sUserLogin) UpdatePasswordRegister(ctx context.Context, token string, password string) (userId int, err error) {
	// 1. token is already verified
	infoOTP, err := s.r.GetInfoOTP(ctx, token)
	if err != nil {
		return response.ErrCodeUserOtpNotExists, err
	}

	// 2. check isVerified OK
	if infoOTP.IsVerified.Int32 != 1 {
		return response.ErrCodeUserOtpNotExists, fmt.Errorf("user otp not verified")
	}

	// 3. check token is exists in user_base

	// 4. add userBase to user_base table
	userSalt, err := crypto.GenarateSalt(16)
	if err != nil {
		return response.ErrCodeUserOtpNotExists, err
	}

	newUserBase, err := s.r.AddUserBase(ctx, database.AddUserBaseParams{
		UserAccount:  infoOTP.VerifyKey,
		UserSalt:     userSalt,
		UserPassword: crypto.HashPassword(password, userSalt),
	})
	if err != nil {
		return response.ErrCodeUserOtpNotExists, err
	}

	newUserId, err := newUserBase.LastInsertId()
	if err != nil {
		return response.ErrCodeUserOtpNotExists, err
	}

	// 5. add newUserId to user_info table
	newUserInfo, err := s.r.AddUserHaveUserId(ctx, database.AddUserHaveUserIdParams{
		UserID:               int32(newUserId),
		UserAccount:          infoOTP.VerifyKey,
		UserNickname:         sql.NullString{String: infoOTP.VerifyKey, Valid: true},
		UserAvatar:           sql.NullString{String: "", Valid: true},
		UserState:            1,
		UserMobile:           sql.NullString{String: "", Valid: true},
		UserGender:           sql.NullInt16{Int16: 0, Valid: true},
		UserBirthday:         sql.NullTime{Time: time.Time{}, Valid: false},
		UserEmail:            sql.NullString{String: infoOTP.VerifyKey, Valid: true},
		UserIsAuthentication: 1,
	})
	if err != nil {
		return response.ErrCodeUserOtpNotExists, err
	}

	newUserId, err = newUserInfo.LastInsertId()
	if err != nil {
		return response.ErrCodeUserOtpNotExists, err
	}

	return int(newUserId), nil
}
