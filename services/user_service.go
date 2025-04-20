package services

import (
	"sample-web/dto"
	"sample-web/models"
	"sample-web/repositories"
)

type UserService interface {
	CreateUser(userRequestDto dto.UserRequest) (dto.UserResponse, error)
	GetUserByEmail(email string) (dto.UserResponse, error)
	UpdateUser(userRequestDto dto.UserRequest) (dto.UserResponse, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (u *userService) CreateUser(userRequestDto dto.UserRequest) (dto.UserResponse, error) {
	user := models.User{
		Email:       userRequestDto.Email,
		Password:    userRequestDto.Password,
		PhoneNumber: userRequestDto.PhoneNumber,
	}
	createdUser, err := u.userRepo.CreateUser(user)
	if err != nil {
		return dto.UserResponse{}, err
	}
	userResponse := dto.UserResponse{
		Id:          createdUser.Id,
		Email:       createdUser.Email,
		PhoneNumber: createdUser.PhoneNumber,
	}
	return userResponse, nil
}

func (u *userService) GetUserByEmail(email string) (dto.UserResponse, error) {
	panic("unimplemented")
}

func (u *userService) UpdateUser(userRequestDto dto.UserRequest) (dto.UserResponse, error) {
	panic("unimplemented")
}
