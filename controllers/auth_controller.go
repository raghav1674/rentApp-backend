package controllers

import (
	"net/http"
	"sample-web/dto"
	"sample-web/errors"
	"sample-web/services"
	"sample-web/utils"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
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

	spanCtx, span := utils.Tracer().Start(ctx.Request.Context(), "AuthController.Register")
	defer span.End()

	var registerRequest dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&registerRequest); err != nil {
		span.RecordError(err)
		ctx.Error(errors.NewAppError(http.StatusBadRequest, "invalid request", err))
		return
	}
	user, err := a.authService.Register(spanCtx, registerRequest)
	if err != nil {
		span.RecordError(err)
		ctx.Error(errors.NewAppError(http.StatusInternalServerError, "email already registered", err))
		return
	}
	span.SetAttributes(
		attribute.String("user_id", user.Id),
	)
	span.AddEvent("User registered")
	ctx.JSON(http.StatusOK, user)
}
