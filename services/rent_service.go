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

type RentService interface {
	CreateRent(ctx context.Context, rentRequest dto.RentRequest) (dto.RentResponse, error)
	GetAllRents(ctx context.Context, userId string, userRole string) (dto.RentResponse, error)
	GetRentById(ctx context.Context, rentId string) (dto.RentResponse, error)
	UpdateRent(ctx context.Context, rentId string,rentRequest dto.RentRequest) (dto.RentResponse, error)
	CloseRent(ctx context.Context, rentId string) (dto.RentResponse, error)
}

type rentService struct {
	rentRepo repositories.RentRepository
	userRepo repositories.UserRepository
}


func NewRentService(rentRepo repositories.RentRepository, userRepo repositories.UserRepository) RentService {
	return &rentService{
		rentRepo: rentRepo,
		userRepo: userRepo,
	}
}

func (r *rentService) CreateRent(ctx context.Context, rentRequest dto.RentRequest) (dto.RentResponse, error) {

	log := utils.GetLogger()
	spanCtx, span := log.Tracer().Start(ctx, "RentService.CreateRent")
	defer span.End()
	
	// get landlord and tenant from the database
	landLord, err := r.userRepo.FindUserById(spanCtx, rentRequest.LandLordId)

	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("Failed to find landlord with %s", err.Error()))
		return dto.RentResponse{}, err
	}

	tenant, err := r.userRepo.FindUserById(spanCtx, rentRequest.TenantId)
	if err != nil {
		log.Error(spanCtx, "Failed to find tenant with %s", err.Error())
		return dto.RentResponse{}, err
	}

	now := time.Now()

	rent := models.Rent{
		LandLord: models.PersonRef{
			Id:   landLord.Id,
			Name: landLord.Name,
		},
		Tenant: models.PersonRef{
			Id:   tenant.Id,
			Name: tenant.Name,
		},
		Location: rentRequest.Location,
		Amount:   rentRequest.Amount,
		Schedule: models.RentSchedule(rentRequest.Schedule),
		Status:   models.RentStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	createdRent, err := r.rentRepo.CreateRent(ctx, rent)

	if err != nil {
		log.Error(spanCtx, "Failed to create rent with %s", err.Error())
		return dto.RentResponse{}, err
	}

	log.Info(spanCtx, "Rent created successfully with ID: %s", createdRent.Id)

	return dto.RentResponse{
		Rents: []models.Rent{createdRent},
	}, nil
}


func (r *rentService) GetAllRents(ctx context.Context, userId string, userRole string) (dto.RentResponse, error) {
	userRoleModel := models.UserRole(userRole)
	rent, err := r.rentRepo.GetAllRents(ctx, userId, userRoleModel)
	if err != nil {
		return dto.RentResponse{}, err
	}
	return dto.RentResponse{
		Rents: rent,
	}, nil
}

// GetRentById implements RentService.
func (r *rentService) GetRentById(ctx context.Context, rentId string) (dto.RentResponse, error) {
	
	log := utils.GetLogger()
	spanCtx, span := log.Tracer().Start(ctx, "RentService.GetRentById")
	defer span.End()

	log.Info(spanCtx, "Rent ID: %s", rentId)
	
	rent, err := r.rentRepo.FindRentById(spanCtx, rentId)
	if err != nil {
		log.Error(spanCtx, "Failed to find rent with %s", err.Error())
		return dto.RentResponse{}, err
	}

	log.Info(spanCtx, "Rent found with ID: %s", rent.Id)

	return dto.RentResponse{
		Rents: []models.Rent{rent},
	}, nil	

}

// UpdateRent implements RentService.
func (r *rentService) UpdateRent(ctx context.Context, rentId string, rentRequest dto.RentRequest) (dto.RentResponse, error) {

	log := utils.GetLogger()
	spanCtx, span := log.Tracer().Start(ctx, "RentService.UpdateRent")
	defer span.End()

	now := time.Now()

	log.Info(spanCtx, "Rent ID: %s", rentId)
	rent, err := r.rentRepo.FindRentById(spanCtx, rentId)
	if err != nil {
		log.Error(spanCtx, "Failed to find rent with %s", err.Error())
		return dto.RentResponse{}, err
	}
	log.Info(spanCtx, "Rent found with ID: %s", rent.Id)

	if rent.Status == models.RentStatusInactive {
		log.Error(spanCtx, "Failed to update rent as it is already closed")
		return dto.RentResponse{}, errors.New("rent is already closed")
	}

	if rentRequest.Location != "" {
		rent.Location = rentRequest.Location
	}
	if rentRequest.Amount != 0 {
		rent.Amount = rentRequest.Amount
	}
	if rentRequest.Schedule != "" {
		rent.Schedule = models.RentSchedule(rentRequest.Schedule)
	}
	if rentRequest.Status != "" {
		rent.Status = models.RentStatus(rentRequest.Status)
	}
	rent.UpdatedAt = now

	updatedRent, err := r.rentRepo.UpdateRent(spanCtx, rent)

	if err != nil {
		log.Error(spanCtx, "Failed to update rent with %s", err.Error())
		return dto.RentResponse{}, err
	}
	
	log.Info(spanCtx, "Rent updated successfully with ID: %s", updatedRent.Id)

	return dto.RentResponse{
		Rents: []models.Rent{updatedRent},
	}, nil

}

// CloseRent implements RentService.
func (r *rentService) CloseRent(ctx context.Context, rentId string) (dto.RentResponse, error) {
	
	rent, err := r.rentRepo.FindRentById(ctx, rentId)
	if err != nil {
		return dto.RentResponse{}, err
	}

	if rent.Status == models.RentStatusInactive {
		return dto.RentResponse{}, errors.New("rent is already closed")
	}

	rent.Status = models.RentStatusInactive

	updatedRent, err := r.rentRepo.UpdateRent(ctx, rent)

	if err != nil {
		return dto.RentResponse{}, err
	}
	return dto.RentResponse{
		Rents: []models.Rent{updatedRent},
	}, nil
}

