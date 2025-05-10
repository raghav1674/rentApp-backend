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

	otp_identifier, err := a.otpService.SendOTP(spanCtx, otpRequest.PhoneNumber)
	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("failed to generate OTP with error %s", err.Error()))
		ctx.Error(customerr.NewAppError(http.StatusInternalServerError, "failed to generate OTP", err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully","identifier": otp_identifier })
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
		ctx.Error(customerr.NewAppError(http.StatusInternalServerError, "failed to verify OTP",err))
		return
	}
	if !isValid {
		log.Error(spanCtx, "invalid OTP")
		ctx.Error(customerr.NewAppError(http.StatusUnauthorized, "invalid OTP", nil))
		return
	}

	log.Info(spanCtx, "OTP verified successfully")

	log.Info(spanCtx,fmt.Sprintf("Checking if user already exists and generate access and refresh tokesn for user with phoneNumber %s",otpRequest.PhoneNumber))

	authResponse, err := a.authService.Login(spanCtx, dto.LoginRequest{
		PhoneNumber: otpRequest.PhoneNumber,
	})

	if err == nil {
		log.Info(spanCtx,"User found")
		ctx.JSON(http.StatusOK, gin.H{
			"otp_verified":    true,
			"user_registered": true,
			"access_token":    authResponse.AccessToken,
			"refresh_token":   authResponse.RefreshToken,
			"error":           false,
		})
		return
	}

	log.Info(spanCtx,fmt.Sprintf("no user found with phone number %s",otpRequest.PhoneNumber))
	
	ctx.JSON(http.StatusOK, gin.H{
		"otp_verified":    true,
		"user_registered": false,
		"error":           err.Error(),
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

	log.Info(spanCtx,fmt.Sprintf("Generating access and refresh tokesn for user %s",user.Id))

	authResponse, err := a.authService.Login(spanCtx, dto.LoginRequest{
		PhoneNumber: user.PhoneNumber,
	})

	if err != nil {
		log.Error(spanCtx,fmt.Sprintf("failed to generate access token for user %s with error %s",user.Id,err.Error()))
		ctx.Error(customerr.NewAppError(http.StatusInternalServerError,"failed to generate access token",err))
		return
	}

	log.Info(spanCtx,fmt.Sprintf("Generated Access/Refresh token for user %s",user.Id))

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":    authResponse.AccessToken,
		"refresh_token":   authResponse.RefreshToken,
		"user": user,
	})
}
