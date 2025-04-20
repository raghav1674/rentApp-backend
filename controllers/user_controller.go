package controllers

import (
	"sample-web/dto"
	"sample-web/services"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	CreateUser(ctx *gin.Context)
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

func (u *userController) CreateUser(ctx *gin.Context) {
	var userRequestDto dto.UserRequest
	if err := ctx.ShouldBindJSON(&userRequestDto); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	createdUser, err := u.userService.CreateUser(userRequestDto)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}
	ctx.JSON(201, createdUser)
}

func (u *userController) GetUserByEmail(ctx *gin.Context) {}

func (u *userController) UpdateUser(ctx *gin.Context) {}
