package dto

import "sample-web/models"

type RentResponse struct {
	Rents []models.Rent `json:"rents"`
}

type RentRequest struct {
	TenantPhoneNumber string  `json:"tenant_phone_number" binding:"required"`
	Title             string  `json:"title" binding:"required"`
	Amount            float64 `json:"amount" binding:"required"`
	Schedule          string  `json:"schedule" binding:"required,oneof=weekly monthly querterly"`
	Status            string  `json:"status"`
	StartDate         string  `json:"start_date" binding:"required"`
	EndDate           string  `json:"end_date" binding:"required"`
}


type RentUpdateRequest struct {
	Title             string  `json:"title" binding:"required"`
	Amount            float64 `json:"amount" binding:"required"`
	Schedule          string  `json:"schedule" binding:"required,oneof=weekly monthly querterly"`
	EndDate           string  `json:"end_date" binding:"required"`
}
