package controllers

import (
	"errors"
	"net/http"
	"sample-web/dto"
	customerr "sample-web/errors"
	"sample-web/services"
	"sample-web/utils"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
)

type AuthController interface {
	GenerateOTP(ctx *gin.Context)
	VerifyOTP(ctx *gin.Context)
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
}

type authController struct {
	authService services.AuthService
	otpService  services.OTPService
}

func NewAuthController(authService services.AuthService, otpService services.OTPService) AuthController {
	return &authController{
		authService: authService,
		otpService:  otpService,
	}
}

func (a *authController) GenerateOTP(ctx *gin.Context) {
	spanCtx, span := utils.Tracer().Start(ctx.Request.Context(), "AuthController.GenerateOTP")
	defer span.End()

	span.AddEvent("GenerateOTPRequestReceived")

	var otpRequest struct {
		PhoneNumber string `json:"phone_number" binding:"required,e164"`
	}
	if err := ctx.ShouldBindJSON(&otpRequest); err != nil {
		span.RecordError(err)
		ctx.Error(customerr.NewAppError(http.StatusBadRequest, customerr.ValidationErrorResponse(err), err))
		return
	}

	_, err := a.otpService.SendOTP(spanCtx, otpRequest.PhoneNumber)
	if err != nil {
		span.RecordError(err)
		ctx.Error(customerr.NewAppError(http.StatusInternalServerError, "failed to generate OTP", err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

func (a *authController) VerifyOTP(ctx *gin.Context) {
	spanCtx, span := utils.Tracer().Start(ctx.Request.Context(), "AuthController.VerifyOTP")
	defer span.End()
	span.AddEvent("VerifyOTPRequestReceived")

	var otpRequest struct {
		PhoneNumber string `json:"phone_number" binding:"required,e164"`
		Code        string `json:"code" binding:"required,min=6,max=6,numeric"`
	}

	if err := ctx.ShouldBindJSON(&otpRequest); err != nil {
		span.RecordError(err)
		span.AddEvent("VerifyOTPRequestFailed")
		ctx.Error(customerr.NewAppError(http.StatusBadRequest, customerr.ValidationErrorResponse(err), err))
		return
	}
	isValid, err := a.otpService.VerifyOTP(spanCtx, otpRequest.PhoneNumber, otpRequest.Code)
	if err != nil {
		span.RecordError(err)
		span.AddEvent("VerifyOTPRequestFailed")
		ctx.Error(customerr.NewAppError(http.StatusInternalServerError, "failed to verify OTP", err))
		return
	}
	if !isValid {
		span.AddEvent("VerifyOTPRequestFailed")
		span.RecordError(errors.New("invalid OTP"))
		ctx.Error(customerr.NewAppError(http.StatusUnauthorized, "invalid OTP", nil))
		return
	}
	span.SetAttributes(
		attribute.String("phone_number", otpRequest.PhoneNumber),
	)
	span.AddEvent("OTP verified successfully")
	ctx.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
}

func (a *authController) Login(ctx *gin.Context) {

	spanCtx, span := utils.Tracer().Start(ctx.Request.Context(), "AuthController.Register")
	defer span.End()

	span.AddEvent("LoginRequestReceived")

	var loginRequest dto.LoginRequest
	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		span.RecordError(err)
		span.AddEvent("LoginRequestFailed")
		ctx.Error(customerr.NewAppError(http.StatusBadRequest, customerr.ValidationErrorResponse(err), err))
		return
	}

	token, err := a.authService.Login(spanCtx, loginRequest)
	if err != nil {
		span.RecordError(err)
		span.AddEvent("LoginRequestFailed")
		ctx.Error(customerr.NewAppError(http.StatusUnauthorized, "invalid credentials", err))
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
		ctx.Error(customerr.NewAppError(http.StatusBadRequest, customerr.ValidationErrorResponse(err), err))
		return
	}
	user, err := a.authService.Register(spanCtx, registerRequest)
	if err != nil {
		span.RecordError(err)
		ctx.Error(customerr.NewAppError(http.StatusBadRequest, "email already registered", err))
		return
	}
	span.SetAttributes(
		attribute.String("user_id", user.Id),
	)
	span.AddEvent("User registered")
	ctx.JSON(http.StatusOK, user)
}
