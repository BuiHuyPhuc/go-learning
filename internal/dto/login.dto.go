package dto

type RegisterRequest struct {
	VerifyKey     string `json:"verify_key"`
	VerifyType    int    `json:"verify_type"`
	VerifyPurpose string `json:"verify_purpose"`
}

type VerifyOTPRequest struct {
	VerifyKey  string `json:"verify_key"`
	VerifyCode string `json:"verify_code"`
}

type VerifyOTPResponse struct {
	Token   string `json:"token"`
	UserId  int    `json:"user_id"`
	Message string `json:"message"`
}

type UpdatePasswordRegisterRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type LoginRequest struct {
	UserAccount  string `json:"user_account"`
	UserPassword string `json:"user_password"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}
