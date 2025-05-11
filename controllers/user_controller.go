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
	GetCurrentUser(ctx *gin.Context)
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

func (u *userController) GetCurrentUser(ctx *gin.Context) {

	log := utils.GetLogger()

	spanCtx, span := log.Tracer().Start(ctx.Request.Context(), "controllers.UserController.GetCurrentUser")
	defer span.End()

	log.Info(spanCtx, "Fetching user_id from context")

	value, exists := ctx.Get("user_id")

	if !exists {
		log.Error(spanCtx, "user_id not found in context")
		ctx.Error(customerr.NewAppError(http.StatusInternalServerError, "user_id not found", nil))
		return
	}

	log.Info(spanCtx, fmt.Sprintf("Retrieving Current User info with user_id %s", value))

	user, err := u.userService.GetUserById(spanCtx, value.(string))
	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("User not found with user_id %s with error %s", user.Id, err.Error()))
		ctx.Error(customerr.NewAppError(http.StatusBadRequest, "user not found", err))
		return
	}
	log.Info(spanCtx, "User Found")
	ctx.JSON(http.StatusOK, user)
}
