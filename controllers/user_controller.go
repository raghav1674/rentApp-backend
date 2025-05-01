package controllers

import (
	"net/http"
	"sample-web/dto"
	customerr "sample-web/errors"
	"sample-web/services"
	"sample-web/utils"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	GetUserByEmail(ctx *gin.Context)
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

// GetUserByEmail implements UserController.
func (u *userController) GetUserByEmail(ctx *gin.Context) {

	spanCtx,span := utils.Tracer().Start(ctx.Request.Context(), "controllers.UserController.GetUserByEmail")
	defer span.End()

	var emailRequest struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := ctx.ShouldBindJSON(&emailRequest); err != nil {
		span.RecordError(err)
		ctx.Error(customerr.NewAppError(http.StatusBadRequest, "Invalid email format",err))
		return
	}

	span.AddEvent("finding user by email")

	user, err := u.userService.GetUserByEmail(spanCtx, emailRequest.Email)
	
	if err != nil {
		span.RecordError(err)
		ctx.Error(customerr.NewAppError(http.StatusInternalServerError, "Error occurred while fetching user information",err))
		return
	}
	ctx.JSON(200, user)
}

// UpdateUser implements UserController.
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
