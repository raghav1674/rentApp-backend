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
	CreateRent(ctx context.Context, landLordId string, rentRequest dto.RentRequest) (dto.RentResponse, error)
	GetAllRents(ctx context.Context, userId string, userRole string) (dto.RentResponse, error)
	GetRentById(ctx context.Context, userId string, rentId string) (dto.RentResponse, error)
	UpdateRent(ctx context.Context, landLordId string, rentId string, rentRequest dto.RentUpdateRequest) (dto.RentResponse, error)
	CloseRent(ctx context.Context, landLordId string,rentId string) (dto.RentResponse, error)
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

func (r *rentService) CreateRent(ctx context.Context, landLordId string, rentRequest dto.RentRequest) (dto.RentResponse, error) {

	log := utils.GetLogger()
	spanCtx, span := log.Tracer().Start(ctx, "RentService.CreateRent")
	defer span.End()

	log.Info(spanCtx, "finding landlord with ID: %s", landLordId)

	// get landlord and tenant from the database
	landLord, err := r.userRepo.FindUserById(spanCtx, landLordId)

	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("Failed to find landlord with %s", err.Error()))
		return dto.RentResponse{}, errors.New("landlord not found")
	}

	log.Info(spanCtx, "landlord found with ID: %s", landLord.Id)

	log.Info(spanCtx, "finding tenant with phone number: %s", rentRequest.TenantPhoneNumber)

	tenant, err := r.userRepo.FindUserByPhoneNumber(spanCtx, rentRequest.TenantPhoneNumber)
	if err != nil {
		log.Error(spanCtx, "Failed to find tenant with %s", err.Error())
		return dto.RentResponse{}, errors.New("tenant not found")
	}

	log.Info(spanCtx, "tenant found with ID: %s", tenant.Id)

	if landLord.Id == tenant.Id {
		log.Error(spanCtx, "landlord and tenant cannot be the same person")
		return dto.RentResponse{}, errors.New("landlord and tenant cannot be the same person")
	}

	now := time.Now()

	startDate,err := time.Parse( "2006-01-02",rentRequest.StartDate)
	if err != nil {
		log.Error(spanCtx, "Failed to parse start date with %s", err.Error())
		return dto.RentResponse{}, errors.New("failed to parse start date")
	}
	endDate,err := time.Parse( "2006-01-02",rentRequest.EndDate)
	if err != nil {
		log.Error(spanCtx, "Failed to parse end date with %s", err.Error())
		return dto.RentResponse{}, errors.New("failed to parse end date")
	}

	if startDate.Before(now) || endDate.Before(startDate) {
		log.Error(spanCtx, "Start date must be after current date and end date must be after start date")
		return dto.RentResponse{}, errors.New("start date must be after current date and end date must be after start date")
	}

	switch models.RentSchedule(rentRequest.Schedule) {
		case models.RentScheduleWeekly:
			if endDate.Sub(startDate).Hours() < 168 || endDate.Sub(startDate).Hours() > 672 {
				log.Error(spanCtx, "Weekly rent must be at least 7 days and at most 28 days")
				return dto.RentResponse{}, errors.New("weekly rent must be at least 7 days and at most 28 days")
			}
		case models.RentScheduleMonthly:
			if endDate.Sub(startDate).Hours() < 672 || endDate.Sub(startDate).Hours() > 2016 {
				log.Error(spanCtx, "Monthly rent must be at least 28 days and at most 84 days")
				return dto.RentResponse{}, errors.New("monthly rent must be at least 28 days and at most 84 days")
			}			
		case models.RentScheduleQuarterly:
			if endDate.Sub(startDate).Hours() < 2016 || endDate.Sub(startDate).Hours() > 6720 {
				log.Error(spanCtx, "Quarterly rent must be at least 84 days and at most 280 days")
				return dto.RentResponse{}, errors.New("quarterly rent must be at least 84 days and at most 280 days")
			}
		default:
			log.Error(spanCtx, "Invalid rent schedule")
			return dto.RentResponse{}, errors.New("invalid rent schedule")
	}


	rent := models.Rent{
		LandLord: models.PersonRef{
			Id:   landLord.Id,
			Name: landLord.Name,
		},
		Tenant: models.PersonRef{
			Id:   tenant.Id,
			Name: tenant.Name,
		},
		Title:     rentRequest.Title,
		Amount:    rentRequest.Amount,
		Schedule:  models.RentSchedule(rentRequest.Schedule),
		Status:    models.RentStatusActive,
		StartDate: startDate,
		EndDate:   endDate,
		CreatedAt: now,
		UpdatedAt: now,
	}

	log.Info(spanCtx, "Creating rent with title: %s", rent.Title)

	createdRent, err := r.rentRepo.CreateRent(ctx, rent)

	if err != nil {
		log.Error(spanCtx, "Failed to create rent with %s", err.Error())
		return dto.RentResponse{}, errors.New("failed to create rent")
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
func (r *rentService) GetRentById(ctx context.Context, userId string, rentId string) (dto.RentResponse, error) {

	log := utils.GetLogger()
	spanCtx, span := log.Tracer().Start(ctx, "RentService.GetRentById")
	defer span.End()

	log.Info(spanCtx, "Rent ID: %s", rentId)

	rent, err := r.rentRepo.FindRentById(spanCtx,userId,rentId)
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
func (r *rentService) UpdateRent(ctx context.Context, landLordId, rentId string, rentRequest dto.RentUpdateRequest) (dto.RentResponse, error) {

	log := utils.GetLogger()
	spanCtx, span := log.Tracer().Start(ctx, "RentService.UpdateRent")
	defer span.End()

	now := time.Now()

	log.Info(spanCtx, "Rent ID: %s", rentId)
	rent, err := r.rentRepo.FindRentById(spanCtx, landLordId, rentId)
	if err != nil {
		log.Error(spanCtx, "Failed to find rent with %s", err.Error())
		return dto.RentResponse{}, err
	}
	log.Info(spanCtx, "Rent found with ID: %s", rent.Id)

	if rent.Status == models.RentStatusInactive {
		log.Error(spanCtx, "Failed to update rent as it is already closed")
		return dto.RentResponse{}, errors.New("rent is already closed")
	}

	if rentRequest.Title != "" {
		rent.Title = rentRequest.Title
	}
	if rentRequest.Amount != 0 {
		rent.Amount = rentRequest.Amount
	}
	if rentRequest.Schedule != "" {
		rent.Schedule = models.RentSchedule(rentRequest.Schedule)
	}
	if rentRequest.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", rentRequest.EndDate)
		if err != nil {
			log.Error(spanCtx, "Failed to parse end date with %s", err.Error())
			return dto.RentResponse{}, errors.New("failed to parse end date")
		}
		if endDate.Before(now) || endDate.Before(rent.StartDate) {
			log.Error(spanCtx, "End date must be after current date and after start date")
			return dto.RentResponse{}, errors.New("end date must be after current date and after start date")
		}
		rent.EndDate = endDate
	}
	rent.UpdatedAt = now

	updatedRent, err := r.rentRepo.UpdateRent(spanCtx, landLordId, rent)

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
func (r *rentService) CloseRent(ctx context.Context, landLordId, rentId string) (dto.RentResponse, error) {

	rent, err := r.rentRepo.FindRentById(ctx, landLordId,rentId)
	if err != nil {
		return dto.RentResponse{}, err
	}

	if rent.Status == models.RentStatusInactive {
		return dto.RentResponse{}, errors.New("rent is already closed")
	}

	rent.Status = models.RentStatusInactive

	updatedRent, err := r.rentRepo.UpdateRent(ctx, landLordId, rent)

	if err != nil {
		return dto.RentResponse{}, err
	}
	return dto.RentResponse{
		Rents: []models.Rent{updatedRent},
	}, nil
}
