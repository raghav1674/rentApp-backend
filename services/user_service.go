package services

import (
	"sample-web/dto"
	"sample-web/mappers"
	"sample-web/repositories"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx *gin.Context, userRequestDto dto.UserRequest) (dto.UserResponse, error)
	GetUserByEmail(ctx *gin.Context, email string) (dto.UserResponse, error)
	UpdateUser(ctx *gin.Context, userRequestDto dto.UserRequest) (dto.UserResponse, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (u *userService) CreateUser(ctx *gin.Context, userRequestDto dto.UserRequest) (dto.UserResponse, error) {
	user := mappers.ToUserModel(userRequestDto)
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequestDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.UserResponse{}, err
	}
	user.Password = string(hasedPassword)
	createdUser, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return dto.UserResponse{}, err
	}
	userResponse := mappers.ToUserResponse(createdUser)
	return userResponse, nil
}

func (u *userService) GetUserByEmail(ctx *gin.Context, email string) (dto.UserResponse, error) {
	user, err := u.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return dto.UserResponse{}, err
	}
	userResponse := mappers.ToUserResponse(user)
	return userResponse, nil
}

func (u *userService) UpdateUser(ctx *gin.Context, userRequestDto dto.UserRequest) (dto.UserResponse, error) {
	user, err := u.userRepo.FindUserByEmail(ctx, userRequestDto.Email)
	if err != nil {
		return dto.UserResponse{}, err
	}
	user.PhoneNumber = userRequestDto.PhoneNumber
	user.Roles = mappers.ToUserRoles(userRequestDto.Roles)
	user.UpdatedAt = time.Now()
	updatedUser, err := u.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return dto.UserResponse{}, err
	}
	userResponse := mappers.ToUserResponse(updatedUser)
	return userResponse, nil
}
