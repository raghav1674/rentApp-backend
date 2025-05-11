package controllers

import (
	"github.com/gin-gonic/gin"
)

type RentRecordController interface {
	CreateRentRecord(ctx *gin.Context)
	GetAllRentRecords(ctx *gin.Context)
	GetRentRecordById(ctx *gin.Context)
	ApproveRentRecord(ctx *gin.Context)
}


