package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController interface {
	GetHealth(ctx *gin.Context)
}

type healthController struct {
}

func NewHealthController() HealthController {
	return &healthController{}
}

func (h *healthController) GetHealth(ctx *gin.Context)  {
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "ok",
	})
}