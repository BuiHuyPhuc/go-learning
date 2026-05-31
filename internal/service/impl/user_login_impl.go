package impl

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-learning/global"
	"go-learning/internal/consts"
	"go-learning/internal/database"
	"go-learning/internal/dto"
	"go-learning/internal/utils"
	"go-learning/internal/utils/auth"
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

func (s *sUserLogin) Login(ctx context.Context, in *dto.LoginRequest) (statusCode int, out dto.LoginResponse, err error) {
	// 1. logic login
	userBase, err := s.r.GetOneUserInfo(ctx, in.UserAccount)
	if err != nil {
		return response.ErrCodeAuthFailed, out, err
	}

	// 2. check password
	if !crypto.MatchPassword(userBase.UserPassword, in.UserPassword, userBase.UserSalt) {
		return response.ErrCodeAuthFailed, out, fmt.Errorf("does not match password")
	}

	// 3. check two-factor authentication
	// isTwoFactorEnable, err := s.r.IsTwoFactorEnabled(ctx, userBase.UserID)
	// if err != nil {
	// 	return response.ErrCodeAuthFailed, out, fmt.Errorf("check two-factor authentication failed: %v", err)
	// }
	// if isTwoFactorEnable > 0 {
	// 	// send otp tp in.TwoFactorEmail
	// 	keyUserLoginTwoFactor := crypto.GetHash("2fa:otp:" + strconv.Itoa(int(userBase.UserID)))
	// 	err = global.Rdb.Set(ctx, keyUserLoginTwoFactor, "111111", time.Duration(consts.TIME_OTP_REGISTER)*time.Minute).Err()
	// 	if err != nil {
	// 		return response.ErrCodeAuthFailed, out, fmt.Errorf("set otp redis failed")
	// 	}

	// 	// sent otp via twofactorEmail
	// 	// get email 2fa
	// 	infoUserTwoFactor, err := s.r.GetTwoFactorMethodByIDAndType(ctx, database.GetTwoFactorMethodByIDAndTypeParams{
	// 		UserID:            userBase.UserID,
	// 		TwoFactorAuthType: 1,
	// 	})
	// 	if err != nil {
	// 		return response.ErrCodeAuthFailed, out, fmt.Errorf("get two-factor method failed: %v", err)
	// 	}
	// 	log.Println("send OTP 2fa to Email::", infoUserTwoFactor.TwoFactorEmail.String)
	// 	go sendto.SendTextEmailOtp([]string{infoUserTwoFactor.TwoFactorEmail.String}, consts.HOST_EMAIL, "111111")

	// 	out.Message = "send OTP 2FA to Email" + infoUserTwoFactor.TwoFactorEmail.String
	// 	return response.ErrCodeSuccess, out, nil
	// }

	// 4. update password login
	go s.r.LoginUserBase(ctx, database.LoginUserBaseParams{
		UserLoginIp:  sql.NullString{String: "127.0.0.1", Valid: true},
		UserAccount:  in.UserAccount,
		UserPassword: in.UserPassword,
	})

	// 5. create UUID user
	subToken := utils.GenerateCliTokenUUID(int(userBase.UserID))
	log.Println("subtoken:", subToken)

	// 6. get user_info table
	infoUser, err := s.r.GetUser(ctx, int32(userBase.UserID))
	if err != nil {
		return response.ErrCodeAuthFailed, out, err
	}
	// convert to json
	infoUserJson, err := json.Marshal(infoUser)
	if err != nil {
		return response.ErrCodeAuthFailed, out, fmt.Errorf("convert to json failed: %v", err)
	}

	// 7. give infoUserJson to redis with key = subToken
	err = global.Rdb.Set(ctx, subToken, infoUserJson, time.Duration(consts.TIME_TOKEN_CACHE)*time.Hour).Err()
	if err != nil {
		return response.ErrCodeAuthFailed, out, err
	}

	// 8. create token
	out.Token, err = auth.CreateToken(subToken)
	if err != nil {
		return response.ErrCodeAuthFailed, out, err
	}

	return response.ErrCodeSuccess, out, nil
}

func (s *sUserLogin) IsTwoFactorEnabled(ctx context.Context, userId int) (statusCode int, out bool, err error) {
	return response.ErrCodeSuccess, true, nil
}

func (s *sUserLogin) SetupTwoFactorAuth(ctx context.Context, in *dto.SetupTwoFactorAuthRequest) (statusCode int, err error) {
	// 1. check IsTwoFactorEnabled -> true -> return
	isTwoFactorEnabled, err := s.r.IsTwoFactorEnabled(ctx, int32(in.UserId))
	if err != nil {
		return response.ErrCodeTwoFactorAuthSetupFailed, err
	}

	if isTwoFactorEnabled > 0 {
		return response.ErrCodeTwoFactorAuthSetupFailed, fmt.Errorf("two-factor authentication is already enabled")
	}

	// 2. create new type Auth
	err = s.r.EnableTwoFactorTypeEmail(ctx, database.EnableTwoFactorTypeEmailParams{
		UserID:            int32(in.UserId),
		TwoFactorAuthType: 1,
		TwoFactorEmail:    sql.NullString{String: in.TwoFactorEmail, Valid: true},
	})
	if err != nil {
		return response.ErrCodeTwoFactorAuthSetupFailed, err
	}

	// 3. send otp to in.TwoFactorEmail
	keyUserTwoFactor := crypto.GetHash("2fa:" + strconv.Itoa(in.UserId))
	go global.Rdb.Set(ctx, keyUserTwoFactor, "123456", time.Duration(consts.TIME_TOKEN_CACHE)*time.Minute).Err()
	// if err != nil {
	// 	return response.ErrCodeTwoFactorAuthSetupFailed, err
	// }

	return response.ErrCodeSuccess, nil
}

func (s *sUserLogin) VerifyTwoFactorAuth(ctx context.Context, in *dto.VerifyTwoFactorAuthRequest) (statusCode int, err error) {
	// 1. check IsTwoFactorEnabled
	isTwoFactorAuth, err := s.r.IsTwoFactorEnabled(ctx, int32(in.UserId))
	if err != nil {
		return response.ErrCodeTwoFactorAuthVerifyFailed, err
	}

	if isTwoFactorAuth > 0 {
		return response.ErrCodeTwoFactorAuthVerifyFailed, fmt.Errorf("two-factor authentication is not enabled")
	}

	// 2. check otp in redis available
	keyUserTwoFactor := crypto.GetHash("2fa:" + strconv.Itoa(in.UserId))
	otpVerifyAuth, err := global.Rdb.Get(ctx, keyUserTwoFactor).Result()
	if err == redis.Nil {
		return response.ErrCodeTwoFactorAuthVerifyFailed, fmt.Errorf("Key %s does not exists", keyUserTwoFactor)
	} else if err != nil {
		return response.ErrCodeTwoFactorAuthVerifyFailed, err
	}

	// 3. check otp
	if otpVerifyAuth != in.TwoFactorCode {
		return response.ErrCodeTwoFactorAuthVerifyFailed, fmt.Errorf("otp does not match")
	}

	// 4. update status
	err = s.r.UpdateTwoFactorStatus(ctx, database.UpdateTwoFactorStatusParams{
		UserID:            int32(in.UserId),
		TwoFactorAuthType: 1,
	})
	if err != nil {
		return response.ErrCodeTwoFactorAuthVerifyFailed, err
	}

	// 5. remove otp
	_, err = global.Rdb.Del(ctx, keyUserTwoFactor).Result()
	if err != nil {
		return response.ErrCodeTwoFactorAuthVerifyFailed, err
	}

	return response.ErrCodeSuccess, nil
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
