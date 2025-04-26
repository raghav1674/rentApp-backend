package dto

type LoginRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required"`
	CurrentRole string `json:"current_role" binding:"required"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
}

type RegisterRequest struct {
	Email       string   `json:"email" binding:"required,email"`
	Password    string   `json:"password" binding:"required,min=6"`
	PhoneNumber string   `json:"phone_number" binding:"required"`
	Roles       []string `json:"roles" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
