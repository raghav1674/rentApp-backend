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

	spanCtx, span := utils.Tracer().Start(ctx.Request.Context(), "AuthController.Register")
	defer span.End()
	
	span.AddEvent("LoginRequestReceived")


	var loginRequest dto.LoginRequest
	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		span.RecordError(err)
		span.AddEvent("LoginRequestFailed")
		ctx.Error(errors.NewAppError(http.StatusBadRequest, "invalid request", err))
		return
	}

	token, err := a.authService.Login(spanCtx, loginRequest)
	if err != nil {
		span.RecordError(err)
		span.AddEvent("LoginRequestFailed")
		ctx.Error(errors.NewAppError(http.StatusUnauthorized, "invalid credentials", err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (a *authController) Register(ctx *gin.Context) {

	spanCtx, span := utils.Tracer().Start(ctx.Request.Context(), "AuthController.Register")
	defer span.End()

	span.AddEvent("RegisterRequestReceived")

	var registerRequest dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&registerRequest); err != nil {
		span.RecordError(err)
		ctx.Error(errors.NewAppError(http.StatusBadRequest, "invalid request", err))
		return
	}
	user, err := a.authService.Register(spanCtx, registerRequest)
	if err != nil {
		span.RecordError(err)
		ctx.Error(errors.NewAppError(http.StatusBadRequest, "email already registered", err))
		return
	}
	span.SetAttributes(
		attribute.String("user_id", user.Id),
	)
	span.AddEvent("User registered")
	ctx.JSON(http.StatusOK, user)
}
