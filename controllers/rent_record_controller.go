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

type RentRecordController interface {
	CreateRentRecord(ctx *gin.Context)
	GetAllRentRecords(ctx *gin.Context)
	GetRentRecordById(ctx *gin.Context)
	ApproveRentRecord(ctx *gin.Context)
	RejectRentRecord(ctx *gin.Context)
}

type rentRecorController struct {
	rentRecordService services.RentRecordService
}

func NewRentRecordController(rentRecordService services.RentRecordService) RentRecordController {
	return &rentRecorController{
		rentRecordService: rentRecordService,
	}
}

func (r *rentRecorController) CreateRentRecord(ctx *gin.Context) {

	log := utils.GetLogger()

	spanCtx, span := log.Tracer().Start(ctx.Request.Context(), "RentRecordController.CreateRentRecord")
	defer span.End()

	rentId := ctx.Param("rent_id")

	tenantId, exists := ctx.Get("user_id")

	if !exists {
		log.Error(spanCtx, "Tenant ID is empty")
		ctx.Error(customerr.NewAppError(http.StatusBadRequest, "Tenant ID is empty", nil))
	}

	var rentRecordRequest dto.RentRecordRequest
	if err := ctx.ShouldBindBodyWithJSON(&rentRecordRequest); err != nil {
		log.Error(spanCtx, fmt.Sprintf("failed to bind request body with %s", err.Error()))
		ctx.Error(customerr.NewAppError(http.StatusBadRequest, "invalid request body", err))
		return
	}

	rentResponse, err := r.rentRecordService.CreateRentRecord(spanCtx, tenantId.(string), rentId, rentRecordRequest)

	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("create rent failed with error %s", err.Error()))
		ctx.Error(customerr.NewAppError(http.StatusInternalServerError, "create rent failed", err))
		return
	}

	log.Info(spanCtx, "rent created successfully")

	ctx.JSON(http.StatusCreated, rentResponse)
}

func (r *rentRecorController) GetAllRentRecords(ctx *gin.Context) {
	
	log := utils.GetLogger()

	rentId := ctx.Param("rent_id")

	spanCtx, span := log.Tracer().Start(ctx.Request.Context(), "RentRecordController.GetAllRentRecords")
	defer span.End()

	userId, exists := ctx.Get("user_id")

	if !exists {
		log.Error(spanCtx, "Tenant ID is empty")
		ctx.Error(customerr.NewAppError(http.StatusBadRequest, "Tenant ID is empty", nil))
		return
	}

	userRole, exists := ctx.Get("current_role")
	if !exists {
		log.Error(spanCtx, "User role is empty")
		ctx.Error(customerr.NewAppError(http.StatusBadRequest, "User role is empty", nil))
		return
	}

	log.Info(spanCtx, fmt.Sprintf("userId: %s, userRole: %s", userId.(string), userRole.(string)))

	rentRecords, err := r.rentRecordService.GetAllRentRecords(spanCtx, userId.(string), userRole.(string),rentId)

	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("get all rent records failed with error %s", err.Error()))
		ctx.Error(customerr.NewAppError(http.StatusInternalServerError, "get all rent records failed", err))
		return
	}

	log.Info(spanCtx, "get all rent records successfully")
	ctx.JSON(http.StatusOK, rentRecords)
	if len(rentRecords) == 0 {
		ctx.JSON(http.StatusNoContent, nil)
		return
	}
	ctx.JSON(http.StatusOK, rentRecords)
}
	
func (r *rentRecorController) GetRentRecordById(ctx *gin.Context) {
	
	log := utils.GetLogger()

	spanCtx, span := log.Tracer().Start(ctx.Request.Context(), "RentRecordController.GetRentRecordById")
	defer span.End()

	rentRecordId := ctx.Param("record_id")

	rentId := ctx.Param("rent_id")

	userId, exists := ctx.Get("user_id")

	if !exists {
		log.Error(spanCtx, "Tenant ID is empty")
		ctx.Error(customerr.NewAppError(http.StatusBadRequest, "Tenant ID is empty", nil))
	}

	rentRecordResponse, err := r.rentRecordService.GetRentRecordById(spanCtx, userId.(string), rentId, rentRecordId)

	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("get rent record failed with error %s", err.Error()))
		ctx.Error(customerr.NewAppError(http.StatusInternalServerError, "get rent record failed", err))
		return
	}

	log.Info(spanCtx, "get rent record successfully")
	ctx.JSON(http.StatusOK, rentRecordResponse)
}

func (r *rentRecorController) ApproveRentRecord(ctx *gin.Context) {
	
	log := utils.GetLogger()

	spanCtx, span := log.Tracer().Start(ctx.Request.Context(), "RentRecordController.ApproveRentRecord")
	defer span.End()

	rentRecordId := ctx.Param("record_id")

	rentId := ctx.Param("rent_id")

	userId, exists := ctx.Get("user_id")

	if !exists {
		log.Error(spanCtx, "Tenant ID is empty")
		ctx.Error(customerr.NewAppError(http.StatusBadRequest, "Tenant ID is empty", nil))
	}

	rentRecordResponse, err := r.rentRecordService.ApproveRentRecord(spanCtx, userId.(string), rentId, rentRecordId)

	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("approve rent record failed with error %s", err.Error()))
		ctx.Error(customerr.NewAppError(http.StatusInternalServerError, "approve rent record failed", err))
		return
	}

	log.Info(spanCtx, "approve rent record successfully")
	ctx.JSON(http.StatusOK, rentRecordResponse)
}

func (r *rentRecorController) RejectRentRecord(ctx *gin.Context) {
	
	log := utils.GetLogger()

	spanCtx, span := log.Tracer().Start(ctx.Request.Context(), "RentRecordController.RejectRentRecord")
	defer span.End()

	rentRecordId := ctx.Param("record_id")

	rentId := ctx.Param("rent_id")

	userId, exists := ctx.Get("user_id")

	if !exists {
		log.Error(spanCtx, "Tenant ID is empty")
		ctx.Error(customerr.NewAppError(http.StatusBadRequest, "Tenant ID is empty", nil))
	}

	rentRecordResponse, err := r.rentRecordService.RejectRentRecord(spanCtx, userId.(string), rentId, rentRecordId)

	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("reject rent record failed with error %s", err.Error()))
		ctx.Error(customerr.NewAppError(http.StatusInternalServerError, "reject rent record failed", err))
		return
	}

	log.Info(spanCtx, "reject rent record successfully")
	ctx.JSON(http.StatusOK, rentRecordResponse)
}
