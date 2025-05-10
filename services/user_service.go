package services

import (
	"context"
	"sample-web/dto"
	"sample-web/mappers"
	"sample-web/repositories"
	"time"
)

type UserService interface {
	CreateUser(ctx context.Context, userRequestDto dto.UserRequest) (dto.UserResponse, error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (dto.UserResponse, error)
	GetUserById(ctx context.Context, userId string) (dto.UserResponse, error)
	UpdateUser(ctx context.Context, userRequestDto dto.UserRequest) (dto.UserResponse, error)
}

type userService struct {
	userRepo repositories.UserRepository
}


func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (u *userService) CreateUser(ctx context.Context, userRequestDto dto.UserRequest) (dto.UserResponse, error) {
	user := mappers.ToUserModel(userRequestDto)
	createdUser, err := u.userRepo.CreateUser(ctx, user)
	if err != nil {
		return dto.UserResponse{}, err
	}
	userResponse := mappers.ToUserResponse(createdUser)
	return userResponse, nil
}


func (u *userService) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (dto.UserResponse, error) {
	user, err := u.userRepo.FindUserByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return dto.UserResponse{}, err
	}
	userResponse := mappers.ToUserResponse(user)
	return userResponse, nil
}


func (u *userService) UpdateUser(ctx context.Context, userRequestDto dto.UserRequest) (dto.UserResponse, error) {
	user, err := u.userRepo.FindUserByPhoneNumber(ctx, userRequestDto.PhoneNumber)
	if err != nil {
		return dto.UserResponse{}, err
	}
	user.PhoneNumber = userRequestDto.PhoneNumber
	user.CurrentRole = mappers.ToUserRole(userRequestDto.CurrentRole)
	user.UpdatedAt = time.Now()
	updatedUser, err := u.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return dto.UserResponse{}, err
	}
	userResponse := mappers.ToUserResponse(updatedUser)
	return userResponse, nil
}

func (u *userService) GetUserById(ctx context.Context, userId string) (dto.UserResponse, error) {
	user, err := u.userRepo.FindUserById(ctx, userId)
	if err != nil {
		return dto.UserResponse{}, err
	}
	userResponse := mappers.ToUserResponse(user)
	return userResponse, nil
}