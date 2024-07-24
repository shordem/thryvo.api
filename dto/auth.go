package dto

type AuthDTO struct {
	Email        string  `json:"email"`
	FirstName    string  `json:"firstname"`
	LastName     string  `json:"lastname"`
	ReferralCode *string `json:"referral_code"`
	Password     string  `json:"password"`
}

type LoginResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
