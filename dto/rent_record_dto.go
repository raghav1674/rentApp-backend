package dto

type RentRecordRequest struct {
	Amount float64 `json:"amount" binding:"required"`
}

type RentRecordResponse struct {
	Id          string  `json:"id"`
	RentId      string  `json:"rent_id"`
	Amount      float64 `json:"amount"`
	SubmittedAt string  `json:"submitted_at"`
	ApprovedAt  string  `json:"approved_at"`
	Status      string  `json:"status"`
}
