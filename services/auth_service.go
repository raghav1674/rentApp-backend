package services

import (
	"context"
	"errors"
	"sample-web/dto"
	"sample-web/mappers"
	"sample-web/models"
	"sample-web/repositories"
	"sample-web/utils"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuthService interface {
	Login(ctx context.Context, loginRequest dto.LoginRequest) (dto.AuthResponse, error)
	Register(ctx context.Context, registerRequest dto.RegisterRequest) (dto.UserResponse, error)
}

type authService struct {
	userRepo   repositories.UserRepository
	jwtService JWTService
}

func NewAuthService(userRepo repositories.UserRepository, jwtSrv JWTService) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtService: jwtSrv,
	}
}

func (a *authService) Login(ctx context.Context, loginRequest dto.LoginRequest) (dto.AuthResponse, error) {

	log := utils.GetLogger()

	spanCtx, span := log.Tracer().Start(ctx, "AuthService.Login")
	defer span.End()

	log.Info(spanCtx, "Finding user by phone number")

	user, err := a.userRepo.FindUserByPhoneNumber(ctx, loginRequest.PhoneNumber)

	if err != nil {
		span.RecordError(err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Info(spanCtx, "User not found")
			return dto.AuthResponse{}, errors.New("user not found")
		} else {
			return dto.AuthResponse{}, err
		}
	}

	claims := CustomClaims{
		UserId:      user.Id.Hex(),
		CurrentRole: loginRequest.CurrentRole,
	}

	authResponse, err := a.jwtService.GenerateToken(ctx, claims)

	if err != nil {
		span.RecordError(err)
		return authResponse, err
	}

	log.Info(spanCtx, "Storing refresh token in database")

	if user.RefreshToken.Token != "" {
		log.Info(spanCtx, "Invalidating existing refresh token")
		user.RefreshToken.IsValid = false
	}

	user.RefreshToken.Token = authResponse.RefreshToken
	user.RefreshToken.IsValid = true
	log.Info(spanCtx, "Updating user with new refresh token")

	user, err = a.userRepo.UpdateUser(ctx, user)

	if err != nil {
		log.Error(spanCtx, err.Error())
		return dto.AuthResponse{}, err
	}

	log.Info(spanCtx, "Refresh token stored successfully")

	log.Info(spanCtx, "Mapping auth response")

	return authResponse, nil
}

func (a *authService) Register(ctx context.Context, registerRequest dto.RegisterRequest) (dto.UserResponse, error) {

	log := utils.GetLogger()

	spanCtx, span := log.Tracer().Start(ctx, "AuthService.Register")
	defer span.End()

	existingUser, err := a.userRepo.FindUserByPhoneNumber(spanCtx, registerRequest.PhoneNumber)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			log.Info(spanCtx, "User not found, proceeding with registration")
		} else {
			log.Error(spanCtx, err.Error())
			return dto.UserResponse{}, err
		}
	}

	if existingUser.PhoneNumber == registerRequest.PhoneNumber {
		return dto.UserResponse{}, errors.New("user already exists")
	}

	now := time.Now()

	user := models.User{
		Name:        registerRequest.Name,
		PhoneNumber: registerRequest.PhoneNumber,
		Roles:       mappers.ToUserRoles(registerRequest.Roles),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	log.Info(spanCtx, "Creating user using repository")

	user, err = a.userRepo.CreateUser(ctx, user)

	if err != nil {
		span.RecordError(err)
		return dto.UserResponse{}, err
	}

	log.Info(spanCtx, "User created successfully")

	log.Info(spanCtx, "Mapping user to response")

	userResponse := mappers.ToUserResponse(user)

	log.Info(spanCtx, "User mapped to response successfully")

	return userResponse, nil
}
