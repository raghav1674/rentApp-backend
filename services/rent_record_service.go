package services

import (
	"context"
	"errors"
	"fmt"
	"sample-web/dto"
	"sample-web/models"
	"sample-web/repositories"
	"sample-web/utils"
	"time"
)

type RentRecordService interface {
	CreateRentRecord(ctx context.Context, tenantId string, rentId string, rentRecordRequest dto.RentRecordRequest) (dto.RentRecordResponse, error)
	GetAllRentRecords(ctx context.Context, userId string,userRole string, rentId string) ([]dto.RentRecordResponse, error)
	GetRentRecordById(ctx context.Context, userId string, rentId string, rentRecordId string) (dto.RentRecordResponse, error)
	ApproveRentRecord(ctx context.Context, landLordId string,rentId string, rentRecordId string) (dto.RentRecordResponse, error)
	RejectRentRecord(ctx context.Context, landLordId string, rentId string,rentRecordId string) (dto.RentRecordResponse, error)
}

type rentRecordService struct {
	rentRecordRepository repositories.RentRecordRepository
	rentRepository       repositories.RentRepository
	userRepository       repositories.UserRepository
}

func NewRentRecordService(rentRecordRepository repositories.RentRecordRepository, rentRepository repositories.RentRepository, userRepository repositories.UserRepository) RentRecordService {
	return &rentRecordService{
		rentRecordRepository: rentRecordRepository,
		rentRepository:       rentRepository,
		userRepository:       userRepository,
	}
}

func (r *rentRecordService) CreateRentRecord(ctx context.Context, tenantId string, rentId string, rentRecordRequest dto.RentRecordRequest) (dto.RentRecordResponse, error) {

	log := utils.GetLogger()

	spanCtx, span := log.Tracer().Start(ctx, "RentRecordService.CreateRentRecord")
	defer span.End()

	if rentId == "" {
		log.Error(spanCtx, "Rent ID is empty")
		return dto.RentRecordResponse{}, errors.New("rent ID is empty")
	}

	_, err := r.userRepository.FindUserById(spanCtx, tenantId)
	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("Error fetching tenant with ID %s: %v", tenantId, err))
		return dto.RentRecordResponse{}, err
	}

	log.Info(spanCtx, fmt.Sprintf("Fetching rent with ID %s", rentId))

	rent, err := r.rentRepository.FindRentById(spanCtx, tenantId, rentId)

	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("Error fetching rent with ID %s: %v", rentId, err))
		return dto.RentRecordResponse{}, err
	}

	log.Info(spanCtx, fmt.Sprintf("Fetched rent with ID %s: %+v", rentId, rent))

	now := time.Now()

	var newRentRecord models.RentRecord
	newRentRecord.Amount = rentRecordRequest.Amount
	newRentRecord.RentId = rent.Id
	newRentRecord.SubmittedAt = now
	newRentRecord.Status = models.RentRecordStatusPending
	newRentRecord.Rent = models.RentInfo{
		Amount:   rent.Amount,
		Schedule: rent.Schedule,
	}
	newRentRecord.LandLord = rent.LandLord
	newRentRecord.Tenant = rent.Tenant
	newRentRecord.CreatedAt = now
	log.Info(spanCtx, fmt.Sprintf("Creating rent record for tenant %s and rent %s", tenantId, rentId))

	rentRecord, err := r.rentRecordRepository.CreateRentRecord(spanCtx, newRentRecord)

	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("Error creating rent record: %v", err))
		return dto.RentRecordResponse{}, err
	}

	log.Info(spanCtx, fmt.Sprintf("Created rent record: %+v", rentRecord))

	return dto.RentRecordResponse{
		Id:          rentRecord.Id.Hex(),
		RentId:      rentRecord.RentId.Hex(),
		Amount:      rentRecord.Amount,
		SubmittedAt: rentRecord.SubmittedAt.Format(time.RFC3339),
		Status:      string(rentRecord.Status),
	}, nil
}

// GetAllRentRecords implements RentRecordService.
func (r *rentRecordService) GetAllRentRecords(ctx context.Context, userId string,userRole string, rentId string) ([]dto.RentRecordResponse, error) {
	
	log := utils.GetLogger()
	spanCtx, span := log.Tracer().Start(ctx, "RentRecordService.GetAllRentRecords")
	defer span.End()
	
	if userId == "" {
		log.Error(spanCtx, "User ID is empty")
		return nil, errors.New("user ID is empty")
	}

	if rentId == "" {
		log.Error(spanCtx, "Rent ID is empty")
		return nil, errors.New("rent ID is empty")
	}

	// check if user and rent exist
	_, err := r.userRepository.FindUserById(spanCtx, userId)
	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("Error fetching user with ID %s: %v", userId, err))
		return nil, err
	}
	log.Info(spanCtx, fmt.Sprintf("Fetching rent with ID %s", rentId))
	rent, err := r.rentRepository.FindRentById(spanCtx, userId, rentId)
	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("Error fetching rent with ID %s: %v", rentId, err))
		return nil, err
	}

	log.Info(spanCtx, fmt.Sprintf("Fetched rent with ID %s: %+v", rentId, rent))

	rentRecords, err := r.rentRecordRepository.GetAllRentRecords(spanCtx, userId, userRole,rentId)
	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("Error fetching rent records: %v", err))
		return nil, err
	}
	log.Info(spanCtx, fmt.Sprintf("Fetched rent records: %+v", rentRecords))
	var rentRecordResponses []dto.RentRecordResponse
	for _, rentRecord := range rentRecords {
		rentRecordResponses = append(rentRecordResponses, dto.RentRecordResponse{
			Id:          rentRecord.Id.Hex(),
			RentId:      rentRecord.RentId.Hex(),
			Amount:      rentRecord.Amount,
			SubmittedAt: rentRecord.SubmittedAt.Format(time.RFC3339),
			ApprovedAt:  rentRecord.ApprovedAt.Format(time.RFC3339),
			Status:      string(rentRecord.Status),
		})
	}
	return rentRecordResponses, nil
}

// GetRentRecordById implements RentRecordService.
func (r *rentRecordService) GetRentRecordById(ctx context.Context, userId string, rentId string,rentRecordId string) (dto.RentRecordResponse, error) {
	
	log := utils.GetLogger()
	spanCtx, span := log.Tracer().Start(ctx, "RentRecordService.GetRentRecordById")
	defer span.End()

	if userId == "" {
		log.Error(spanCtx, "User ID is empty")
		return dto.RentRecordResponse{}, errors.New("user ID is empty")
	}

	if rentRecordId == "" {
		log.Error(spanCtx, "Rent Record ID is empty")
		return dto.RentRecordResponse{}, errors.New("rent record ID is empty")
	}

	_, err := r.userRepository.FindUserById(spanCtx, userId)
	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("Error fetching user with ID %s: %v", userId, err))
		return dto.RentRecordResponse{}, err
	}

	rent, err := r.rentRepository.FindRentById(spanCtx, userId, rentId)
	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("Error fetching rent with ID %s: %v", rentId, err))
		return dto.RentRecordResponse{}, err
	}

	log.Info(spanCtx, fmt.Sprintf("Fetched rent with ID %s: %+v", rentId, rent))

	rentRecord, err := r.rentRecordRepository.GetRentRecordById(spanCtx, rentRecordId)
	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("Error fetching rent record with ID %s: %v", rentRecordId, err))
		return dto.RentRecordResponse{}, err
	}

	log.Info(spanCtx, fmt.Sprintf("Fetched rent record with ID %s: %+v", rentRecordId, rentRecord))

	return dto.RentRecordResponse{
		Id:          rentRecord.Id.Hex(),
		RentId:      rentRecord.RentId.Hex(),
		Amount:      rentRecord.Amount,
		SubmittedAt: rentRecord.SubmittedAt.Format(time.RFC3339),
		Status:      string(rentRecord.Status),
	}, nil
}

// ApproveRentRecord implements RentRecordService.
func (r *rentRecordService) ApproveRentRecord(ctx context.Context, landLordId string, rentId string,rentRecordId string) (dto.RentRecordResponse, error) {
	
	log := utils.GetLogger()
	spanCtx, span := log.Tracer().Start(ctx, "RentRecordService.ApproveRentRecord")
	defer span.End()

	if landLordId == "" {
		log.Error(spanCtx, "Landlord ID is empty")
		return dto.RentRecordResponse{}, errors.New("landlord ID is empty")
	}

	if rentRecordId == "" {
		log.Error(spanCtx, "Rent Record ID is empty")
		return dto.RentRecordResponse{}, errors.New("rent record ID is empty")
	}

	_, err := r.userRepository.FindUserById(spanCtx, landLordId)
	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("Error fetching landlord with ID %s: %v", landLordId, err))
		return dto.RentRecordResponse{}, err
	}

	rent, err := r.rentRepository.FindRentById(spanCtx, landLordId, rentId)
	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("Error fetching rent with ID %s: %v", rentId, err))
		return dto.RentRecordResponse{}, err
	}

	log.Info(spanCtx, fmt.Sprintf("Fetched rent with ID %s: %+v", rentId, rent))
	rentRecord, err := r.rentRecordRepository.GetRentRecordById(spanCtx, rentRecordId)
	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("Error fetching rent record with ID %s: %v", rentRecordId, err))
		return dto.RentRecordResponse{}, err
	}

	log.Info(spanCtx, fmt.Sprintf("Fetched rent record with ID %s: %+v", rentRecordId, rentRecord))
	if rentRecord.Status != models.RentRecordStatusPending {
		log.Error(spanCtx, fmt.Sprintf("Rent record with ID %s is not pending", rentRecordId))
		return dto.RentRecordResponse{}, errors.New("rent record is not pending")
	}
	now := time.Now()
	rentRecord.Status = models.RentRecordStatusApproved
	rentRecord.UpdatedAt = now
	rentRecord.ApprovedAt = now

	
	updatedRentRecord, err := r.rentRecordRepository.UpdateRentRecord(spanCtx, landLordId, rentRecordId,rentRecord)

	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("Error updating rent record with ID %s: %v", rentRecordId, err))
		return dto.RentRecordResponse{}, err
	}

	log.Info(spanCtx, fmt.Sprintf("Updated rent record with ID %s: %+v", rentRecordId, updatedRentRecord))
	return dto.RentRecordResponse{
		Id:          updatedRentRecord.Id.Hex(),
		RentId:      updatedRentRecord.RentId.Hex(),
		Amount:      updatedRentRecord.Amount,
		SubmittedAt: updatedRentRecord.SubmittedAt.Format(time.RFC3339),
		ApprovedAt:  updatedRentRecord.ApprovedAt.Format(time.RFC3339),
		Status:      string(updatedRentRecord.Status),
	}, nil
}

// RejectRentRecord implements RentRecordService.
func (r *rentRecordService) RejectRentRecord(ctx context.Context, landLordId string, rentId string,rentRecordId string) (dto.RentRecordResponse, error) {
	
	log := utils.GetLogger()
	spanCtx, span := log.Tracer().Start(ctx, "RentRecordService.RejectRentRecord")
	defer span.End()

	if landLordId == "" {
		log.Error(spanCtx, "Landlord ID is empty")
		return dto.RentRecordResponse{}, errors.New("landlord ID is empty")
	}

	if rentRecordId == "" {
		log.Error(spanCtx, "Rent Record ID is empty")
		return dto.RentRecordResponse{}, errors.New("rent record ID is empty")
	}

	_, err := r.userRepository.FindUserById(spanCtx, landLordId)
	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("Error fetching landlord with ID %s: %v", landLordId, err))
		return dto.RentRecordResponse{}, err
	}

	rent, err := r.rentRepository.FindRentById(spanCtx, landLordId, rentId)
	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("Error fetching rent with ID %s: %v", rentId, err))
		return dto.RentRecordResponse{}, err
	}

	log.Info(spanCtx, fmt.Sprintf("Fetched rent with ID %s: %+v", rentId, rent))
	rentRecord, err := r.rentRecordRepository.GetRentRecordById(spanCtx, rentRecordId)
	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("Error fetching rent record with ID %s: %v", rentRecordId, err))
		return dto.RentRecordResponse{}, err
	}

	log.Info(spanCtx, fmt.Sprintf("Fetched rent record with ID %s: %+v", rentRecordId, rentRecord))
	if rentRecord.Status != models.RentRecordStatusPending {
		log.Error(spanCtx, fmt.Sprintf("Rent record with ID %s is not pending", rentRecordId))
		return dto.RentRecordResponse{}, errors.New("rent record is not pending")
	}

	now := time.Now()
	rentRecord.Status = models.RentRecordStatusRejected
	rentRecord.UpdatedAt = now
	// rentRecord.ApprovedAt = time.Now()

	updatedRentRecord, err := r.rentRecordRepository.UpdateRentRecord(spanCtx, landLordId,rentRecordId,rentRecord)

	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("Error updating rent record with ID %s: %v", rentRecordId, err))
		return dto.RentRecordResponse{}, err
	}
	log.Info(spanCtx, fmt.Sprintf("Updated rent record with ID %s: %+v", rentRecordId, updatedRentRecord))

	return dto.RentRecordResponse{
		Id:          updatedRentRecord.Id.Hex(),
		RentId:      updatedRentRecord.RentId.Hex(),
		Amount:      updatedRentRecord.Amount,
		SubmittedAt: updatedRentRecord.SubmittedAt.Format(time.RFC3339),
		Status:      string(updatedRentRecord.Status),
	}, nil
}
