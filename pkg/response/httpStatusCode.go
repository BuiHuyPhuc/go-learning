package response

const (
	ErrCodeSuccess      = 20001 // success
	ErrCodeParamInvalid = 20003 // param is invalid
	ErrInvalidToken     = 30001 // token is invalid
	ErrInvalidOTP       = 30002 // otp is invalid
	ErrSendEmailOTP     = 30003 // failed to send mail OTP

	// Authentication
	ErrCodeAuthFailed = 40005

	// Register
	ErrCodeUserHasExists = 50001 // user has already exists

	// Login
	ErrCodeOtpNotExists     = 60009
	ErrCodeUserOtpNotExists = 60008
)

var msg = map[int]string{
	ErrCodeSuccess:      "success",
	ErrCodeParamInvalid: "param is invalid",
	ErrInvalidToken:     "token is invalid",
	ErrInvalidOTP:       "otp is invalid",
	ErrSendEmailOTP:     "failed to send mail OTP",

	ErrCodeUserHasExists: "user has already exists",

	ErrCodeOtpNotExists:     "otp exists but not registered",
	ErrCodeUserOtpNotExists: "user otp not exists",

	ErrCodeAuthFailed: "authentication failed",
}
