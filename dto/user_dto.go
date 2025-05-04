package dto

type UserRequest struct {
	PhoneNumber string   `json:"phone_number" binding:"required"`
	Roles       []string `json:"roles" binding:"required"`
}

type UserResponse struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	PhoneNumber string   `json:"phone_number"`
	Roles       []string `json:"roles"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}
