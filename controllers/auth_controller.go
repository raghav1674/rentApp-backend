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

type AuthController interface {
	GenerateOTP(ctx *gin.Context)
	VerifyOTP(ctx *gin.Context)
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

	log := utils.GetLogger()

	spanCtx, span := log.Tracer().Start(ctx.Request.Context(), "AuthController.GenerateOTP")
	defer span.End()

	log.Info(spanCtx, "Generate OTP Request Received")

	var otpRequest struct {
		PhoneNumber string `json:"phone_number" binding:"required,e164"`
	}
	if err := ctx.ShouldBindJSON(&otpRequest); err != nil {
		log.Error(spanCtx, fmt.Sprintf("invalid phone number with error %s", err.Error()))
		ctx.Error(customerr.NewAppError(http.StatusBadRequest, "invalid phone number", err))
		return
	}

	_, err := a.otpService.SendOTP(spanCtx, otpRequest.PhoneNumber)
	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("failed to generate OTP with error %s", err.Error()))
		ctx.Error(customerr.NewAppError(http.StatusInternalServerError, "failed to generate OTP", err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

func (a *authController) VerifyOTP(ctx *gin.Context) {

	log := utils.GetLogger()
	spanCtx, span := log.Tracer().Start(ctx.Request.Context(), "AuthController.VerifyOTP")
	defer span.End()

	log.Info(spanCtx, "Verify OTP Request Received")

	var otpRequest struct {
		PhoneNumber string `json:"phone_number" binding:"required,e164"`
		Code        string `json:"code" binding:"required,min=6,max=6,numeric"`
	}

	if err := ctx.ShouldBindJSON(&otpRequest); err != nil {
		log.Error(spanCtx, fmt.Sprintf("invalid phone number or otp code with error %s", err.Error()))
		ctx.Error(customerr.NewAppError(http.StatusBadRequest, "invalid phone number or otp code", err))
		return
	}

	isValid, err := a.otpService.VerifyOTP(spanCtx, otpRequest.PhoneNumber, otpRequest.Code)

	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("failed to verify OTP with error %s", err.Error()))
		ctx.Error(customerr.NewAppError(http.StatusInternalServerError, "failed to verify OTP", err))
		return
	}
	if !isValid {
		log.Error(spanCtx, "invalid OTP")
		ctx.Error(customerr.NewAppError(http.StatusUnauthorized, "invalid OTP", nil))
		return
	}

	log.Info(spanCtx, "OTP verified successfully")

	authResponse, err := a.authService.Login(spanCtx, dto.LoginRequest{
		PhoneNumber: otpRequest.PhoneNumber,
	})

	if err == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"otp_verified":    true,
			"user_registered": true,
			"access_token":    authResponse.AccessToken,
			"refresh_token":   authResponse.RefreshToken,
			"error":           false,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"otp_verified":    true,
		"user_registered": false,
		"error":           false,
	})
}

func (a *authController) Register(ctx *gin.Context) {

	log := utils.GetLogger()

	spanCtx, span := log.Tracer().Start(ctx.Request.Context(), "AuthController.Register")
	defer span.End()

	log.Info(spanCtx, "Register Request Received")

	var registerRequest dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&registerRequest); err != nil {
		log.Error(spanCtx, fmt.Sprintf("invalid request with error %s", err.Error()))
		ctx.Error(customerr.NewAppError(http.StatusBadRequest, "invalid request", err))
		return
	}
	user, err := a.authService.Register(spanCtx, registerRequest)
	if err != nil {
		log.Error(spanCtx, "register request failed with error %s", err.Error())
		ctx.Error(customerr.NewAppError(http.StatusBadRequest, "register request failed", err))
		return
	}

	log.Info(spanCtx, "User registered")
	ctx.JSON(http.StatusOK, user)
}
