package controllers

import (
	"fmt"
	"sample-web/dto"
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

	spanCtx, span := log.Tracer().Start(ctx, "RentController.CreateRent")
	defer span.End()

	// Bind the request body to the RentRequest struct
	var rentRequest dto.RentRequest
	if err := ctx.ShouldBindJSON(&rentRequest); err != nil {
		log.Error(spanCtx, fmt.Sprintf("Failed to bind request body with %s",err.Error()))
		ctx.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}
	// Call the service to create a new rent
	rent, err := r.rentService.CreateRent(ctx, rentRequest)
	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("Failed to create rent with %s",err.Error()))
		ctx.JSON(500, gin.H{"error": "Failed to create rent"})
		return
	}
	ctx.JSON(201, rent)
}

// GetAllRents implements RentController.
func (r *rentController) GetAllRents(ctx *gin.Context) {
	// Get the user ID from the context
	userId, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(400, gin.H{"error": "User ID not found in context"})
		return
	}
	// Get the user role from the context
	userRole, exists := ctx.Get("current_role")
	if !exists {
		ctx.JSON(400, gin.H{"error": "User role not found in context"})
		return
	}
	// Call the service to get all rents
	rents, err := r.rentService.GetAllRents(ctx, userId.(string), userRole.(string))
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to get rents"})
		return
	}

	ctx.JSON(200, rents)
}

// GetRentById implements RentController.
func (r *rentController) GetRentById(ctx *gin.Context) {
	// Get the rent ID from the URL parameters
	rentId := ctx.Param("rent_id")
	if rentId == "" {
		ctx.JSON(400, gin.H{"error": "Rent ID is required"})
		return
	}

	// Call the service to get the rent by ID
	rent, err := r.rentService.GetRentById(ctx, rentId)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to get rent"})
		return
	}

	ctx.JSON(200, rent)
}

// UpdateRent implements RentController.
func (r *rentController) UpdateRent(ctx *gin.Context) {
	// Get the rent ID from the URL parameters
	rentId := ctx.Param("rent_id")
	if rentId == "" {
		ctx.JSON(400, gin.H{"error": "Rent ID is required"})
		return
	}

	// Bind the request body to the RentRequest struct
	var rentRequest dto.RentRequest
	if err := ctx.ShouldBindJSON(&rentRequest); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	// Call the service to update the rent
	rent, err := r.rentService.UpdateRent(ctx, rentId, rentRequest)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to update rent"})
		return
	}

	ctx.JSON(200, rent)
}


// CloseRent implements RentController.
func (r *rentController) CloseRent(ctx *gin.Context) {
	// Get the rent ID from the URL parameters
	rentId := ctx.Param("rent_id")
	if rentId == "" {
		ctx.JSON(400, gin.H{"error": "Rent ID is required"})
		return
	}

	// Call the service to close the rent
	rent, err := r.rentService.CloseRent(ctx, rentId)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to close rent"})
		return
	}

	ctx.JSON(200, rent)
}

// SummariseRent implements RentController.
func (r *rentController) SummariseRent(ctx *gin.Context) {
	panic("unimplemented")
}

