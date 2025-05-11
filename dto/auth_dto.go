package dto

type LoginRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required,e164"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterRequest struct {
	Name        string `json:"name" binding:"required,min=1"`
	PhoneNumber string `json:"phone_number" binding:"required,e164"`
	CurrentRole string `json:"current_role" binding:"required,oneof=landlord tenant"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
