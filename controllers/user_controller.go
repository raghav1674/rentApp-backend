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

type UserController interface {
	GetUserByPhoneNumber(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
}

type userController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return &userController{
		userService: userService,
	}
}

func (u *userController) GetUserByPhoneNumber(ctx *gin.Context) {

	log := utils.GetLogger()

	spanCtx, span := log.Tracer().Start(ctx.Request.Context(), "controllers.UserController.GetUserByEmail")
	defer span.End()

	var phoneNumberRequest struct {
		PhoneNumber string `json:"phone_number" binding:"required,e164"`
	}
	if err := ctx.ShouldBindJSON(&phoneNumberRequest); err != nil {
		log.Error(spanCtx, err.Error())
		ctx.Error(customerr.NewAppError(http.StatusBadRequest, "Invalid phone format", err))
		return
	}

	log.Info(spanCtx, fmt.Sprintf("finding user by phone number %s", phoneNumberRequest.PhoneNumber))

	user, err := u.userService.GetUserByPhoneNumber(spanCtx, phoneNumberRequest.PhoneNumber)

	if err != nil {
		log.Error(spanCtx, err.Error())
		ctx.Error(customerr.NewAppError(http.StatusInternalServerError, "Error occurred while fetching user information", err))
		return
	}
	ctx.JSON(200, user)
}

func (u *userController) UpdateUser(ctx *gin.Context) {
	var userRequestDto dto.UserRequest
	if err := ctx.ShouldBindJSON(&userRequestDto); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user, err := u.userService.UpdateUser(ctx, userRequestDto)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, user)
}
