package response

const (
	ErrCodeSuccess      = 20001 // success
	ErrCodeParamInvalid = 20003 // email is invalid
	ErrInvalidToken     = 30001 // token is invalid

	// Register
	ErrCodeUserHasExists = 50001 // user has already exists
)

var msg = map[int]string{
	ErrCodeSuccess:      "success",
	ErrCodeParamInvalid: "email is invalid",
	ErrInvalidToken:     "token is invalid",

	ErrCodeUserHasExists: "user has already exists",
}
