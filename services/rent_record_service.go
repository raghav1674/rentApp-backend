package services

import (
	"context"
	"sample-web/dto"
)

type RentRecordService interface {
	CreateRentRecord(ctx context.Context, tenantId string, rentId string, rentRecordRequest dto.RentRecordRequest) (dto.RentRecordResponse, error)
	GetAllRentRecords(ctx context.Context, userId string, rentId string, userRole string) ([]dto.RentRecordResponse, error)
	GetRentRecordById(ctx context.Context, userId string, rentRecordId string) (dto.RentRecordResponse, error)
	ApproveRentRecord(ctx context.Context, landLordId string, rentRecordId string) (dto.RentRecordResponse, error)
}


