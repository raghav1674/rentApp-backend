package dto

import "sample-web/models"

type RentRecordRequest struct {
	RentId string `json:"rent_id" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}

type RentRecordResponse struct {
	Id          string  `json:"id"`
	RentId      string  `json:"rent_id"`
	Rent        models.RentInfo `json:"rent"`
	Amount      float64 `json:"amount"`
	SubmittedAt string  `json:"submitted_at"`
	ApprovedAt  string  `json:"approved_at"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

