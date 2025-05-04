package dto

import "sample-web/models"


type RentResponse struct {
	Rents []models.Rent `json:"rents"`
}

type RentRequest struct {
	LandLordId string `json:"landlord_id" binding:"required"`
	TenantId   string `json:"tenant_id" binding:"required"`
	Location   string `json:"location" binding:"required"`
	Amount     float64 `json:"amount" binding:"required"`
	Schedule   string `json:"schedule" binding:"required,oneof=monthly bimonthly"`
	Status     string `json:"status"`
}