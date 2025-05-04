package services

import (
	"context"
	"errors"
	"sample-web/dto"
	"sample-web/models"
	"sample-web/repositories"
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

	// get landlord and tenant from the database
	landLord, err := r.userRepo.FindUserById(ctx, rentRequest.LandLordId)

	if err != nil {
		return dto.RentResponse{}, err
	}

	tenant, err := r.userRepo.FindUserById(ctx, rentRequest.TenantId)
	if err != nil {
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
		return dto.RentResponse{}, err
	}

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
	rent, err := r.rentRepo.FindRentById(ctx, rentId)
	if err != nil {
		return dto.RentResponse{}, err
	}
	return dto.RentResponse{
		Rents: []models.Rent{rent},
	}, nil	

}

// UpdateRent implements RentService.
func (r *rentService) UpdateRent(ctx context.Context, rentId string, rentRequest dto.RentRequest) (dto.RentResponse, error) {

	now := time.Now()

	rent, err := r.rentRepo.FindRentById(ctx, rentId)
	if err != nil {
		return dto.RentResponse{}, err
	}

	if rent.Status == models.RentStatusInactive {
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

	updatedRent, err := r.rentRepo.UpdateRent(ctx, rent)

	if err != nil {
		return dto.RentResponse{}, err
	}

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

