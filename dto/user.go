package dto

type UserDTO struct {
	DTO

	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	IsEmailVerified bool   `json:"is_email_verified"`
	Password        string `json:"password"`
	Role            string `json:"role"`
}

type VerificationCodeDTO struct {
	DTO

	Code   string  `json:"code"`
	UserID string  `json:"user_id"`
	User   UserDTO `json:"user"`
}
