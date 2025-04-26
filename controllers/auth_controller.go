package controllers

import (
	"net/http"
	"sample-web/dto"
	"sample-web/services"

	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}
type authController struct {
	authService services.AuthService
}

func NewAuthController(authService services.AuthService) AuthController {
	return &authController{
		authService: authService,
	}
}

func (a *authController) Login(ctx *gin.Context) {
	var loginRequest dto.LoginRequest
	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := a.authService.Login(ctx, loginRequest)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (a *authController) Register(ctx *gin.Context) {
	var registerRequest dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&registerRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := a.authService.Register(ctx, registerRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}
