package controllers

import (
	"fmt"
	"net/http"
	"sample-web/dto"
	customerr "sample-web/errors"
	"sample-web/services"
	"sample-web/utils"

	"github.com/gin-gonic/gin"
)

type RentController interface {
	CreateRent(ctx *gin.Context)
	GetAllRents(ctx *gin.Context)
	GetRentById(ctx *gin.Context)
	UpdateRent(ctx *gin.Context)
	CloseRent(ctx *gin.Context)
	SummariseRent(ctx *gin.Context)
}

type rentController struct {
	rentService services.RentService
}


func NewRentController(rentService services.RentService) RentController {
	return &rentController{
		rentService: rentService,
	}
}

func (r *rentController) CreateRent(ctx *gin.Context) {

	log := utils.GetLogger()

	spanCtx, span := log.Tracer().Start(ctx.Request.Context(), "RentController.CreateRent")
	defer span.End()

	// Bind the request body to the RentRequest struct
	var rentRequest dto.RentRequest
	if err := ctx.ShouldBindJSON(&rentRequest); err != nil {
		log.Error(spanCtx, fmt.Sprintf("Failed to bind request body with %s",err.Error()))
		ctx.Error(customerr.NewAppError(http.StatusBadRequest, "Invalid request body", err))
		return
	}
	// Call the service to create a new rent
	rent, err := r.rentService.CreateRent(spanCtx, rentRequest)
	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("Failed to create rent with %s",err.Error()))
		ctx.Error(customerr.NewAppError(http.StatusInternalServerError, "Failed to create rent", err))
		return
	}

	log.Info(spanCtx, fmt.Sprintf("Rent created successfully with ID: %s", rent.Rents[0].Id))
	ctx.JSON(201, rent)
}

// GetAllRents implements RentController.
func (r *rentController) GetAllRents(ctx *gin.Context) {

	log := utils.GetLogger()

	spanCtx, span := log.Tracer().Start(ctx.Request.Context(), "")
	defer span.End()

	// Get the user ID from the context
	userId, exists := ctx.Get("user_id")
	if !exists {
		log.Error(spanCtx, "User ID not found in context")
		ctx.Error(customerr.NewAppError(http.StatusBadRequest, "User ID not found in context", nil))
		return
	}
	// Get the user role from the context
	userRole, exists := ctx.Get("current_role")
	if !exists {
		log.Error(spanCtx, "User role not found in context")
		ctx.Error(customerr.NewAppError(http.StatusBadRequest, "User role not found in context", nil))
		return
	}
	// Call the service to get all rents
	rents, err := r.rentService.GetAllRents(spanCtx, userId.(string), userRole.(string))
	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("Failed to get rents with %s",err.Error()))
		ctx.Error(customerr.NewAppError(http.StatusInternalServerError, "Failed to get rents", err))
		return
	}

	log.Info(spanCtx, fmt.Sprintf("Rents retrieved successfully for user ID: %s", userId))

	ctx.JSON(200, rents)
}

// GetRentById implements RentController.
func (r *rentController) GetRentById(ctx *gin.Context) {

	log := utils.GetLogger()

	spanCtx, span := log.Tracer().Start(ctx.Request.Context(), "RentController.GetRentById")
	defer span.End()
	// Get the rent ID from the URL parameters
	rentId := ctx.Param("rent_id")
	if rentId == "" {
		log.Error(spanCtx, "Rent ID not provided")
		ctx.Error(customerr.NewAppError(http.StatusBadRequest, "Rent ID not provided", nil))
		return
	}

	// Call the service to get the rent by ID
	rent, err := r.rentService.GetRentById(spanCtx, rentId)
	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("Failed to get rent with %s",err.Error()))
		ctx.Error(customerr.NewAppError(http.StatusInternalServerError, "Failed to get rent", err))
		return
	}

	log.Info(spanCtx, fmt.Sprintf("Rent retrieved successfully with ID: %s", rentId))

	ctx.JSON(200, rent)
}

// UpdateRent implements RentController.
func (r *rentController) UpdateRent(ctx *gin.Context) {

	log := utils.GetLogger()

	spanCtx, span := log.Tracer().Start(ctx.Request.Context(), "RentController.UpdateRent")
	defer span.End()

	// Get the rent ID from the URL parameters
	rentId := ctx.Param("rent_id")
	if rentId == "" {
		log.Error(spanCtx, "Rent ID not provided")
		ctx.Error(customerr.NewAppError(http.StatusBadRequest, "Rent ID not provided", nil))
		return
	}

	// Bind the request body to the RentRequest struct
	var rentRequest dto.RentRequest
	if err := ctx.ShouldBindJSON(&rentRequest); err != nil {
		log.Error(spanCtx, fmt.Sprintf("Failed to bind request body with %s",err.Error()))
		ctx.Error(customerr.NewAppError(http.StatusBadRequest, "Invalid request body", err))
		return
	}

	// Call the service to update the rent
	rent, err := r.rentService.UpdateRent(spanCtx, rentId, rentRequest)
	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("Failed to update rent with %s",err.Error()))
		ctx.Error(customerr.NewAppError(http.StatusInternalServerError, "Failed to update rent", err))
		return
	}

	log.Info(spanCtx, fmt.Sprintf("Rent updated successfully with ID: %s", rentId))

	ctx.JSON(200, rent)
}


// CloseRent implements RentController.
func (r *rentController) CloseRent(ctx *gin.Context) {

	log := utils.GetLogger()

	spanCtx, span := log.Tracer().Start(ctx.Request.Context(), "RentController.CloseRent")
	defer span.End()


	// Get the rent ID from the URL parameters
	rentId := ctx.Param("rent_id")
	if rentId == "" {
		log.Error(spanCtx, "Rent ID is required")
		ctx.Error(customerr.NewAppError(http.StatusBadRequest, "Rent ID is required", nil))
		return
	}

	// Call the service to close the rent
	rent, err := r.rentService.CloseRent(spanCtx, rentId)
	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("Failed to close rent with %s",err.Error()))
		ctx.Error(customerr.NewAppError(http.StatusInternalServerError, "Failed to close rent", err))
		return
	}
	log.Info(spanCtx, fmt.Sprintf("Rent closed successfully with ID: %s", rentId))
	ctx.JSON(200, rent)
}

// SummariseRent implements RentController.
func (r *rentController) SummariseRent(ctx *gin.Context) {
	panic("unimplemented")
}

