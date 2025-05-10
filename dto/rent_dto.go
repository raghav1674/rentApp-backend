package dto

import "sample-web/models"


type RentResponse struct {
	Rents []models.Rent `json:"rents"`
}

type RentRequest struct {
	LandLordId string `json:"landlord_id" binding:"required"`
	TenantId   string `json:"tenant_id" binding:"required"`
	Title      string `json:"title" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
	Schedule   string `json:"schedule" binding:"required,oneof=monthly bimonthly"`
	Status     string `json:"status"`
	StartDate  string `json:"start_date" binding:"required"`
	EndDate    string `json:"end_date" binding:"required"`
}