package dto

type UserRequest struct {
	Email       string   `json:"email" binding:"required,email"`
	Password    string   `json:"password" binding:"required,min=6"`
	PhoneNumber string   `json:"phone_number" binding:"required"`
	Roles       []string `json:"roles" binding:"required"`
}

type UserResponse struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	PhoneNumber string   `json:"phone_number"`
	Roles       []string `json:"roles"`
	CreatedAt   string   `json:"created_at"`
	UpdatedAt   string   `json:"updated_at"`
}
