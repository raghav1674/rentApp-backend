package dto

type UserRequest struct {
	PhoneNumber string   `json:"phone_number" binding:"required"`
	CurrentRole string   `json:"current_role" binding:"required"`
}

type UserResponse struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	PhoneNumber string   `json:"phone_number"`
}
