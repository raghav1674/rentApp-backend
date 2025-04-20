package dto

type UserRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
}

type UserResponse struct {
	Id          string `json:"user_id"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}
