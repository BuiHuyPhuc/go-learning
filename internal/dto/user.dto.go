package dto

type UserRegistratorRequest struct {
	Email   string `json:"email" binding:"required,email"`
	Purpose string `json:"purpose" binding:"required"`
}
