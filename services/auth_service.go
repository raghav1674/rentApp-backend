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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
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

	ctx, span := utils.Tracer().Start(ctx, "AuthService.Login")
	defer span.End()

	span.AddEvent("Finding user by email", trace.WithAttributes(
		attribute.String("email", loginRequest.Email)),
	)

	user, err := a.userRepo.FindUserByEmail(ctx, loginRequest.Email)

	if err != nil {
		span.RecordError(err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			span.AddEvent("User not found")
			return dto.AuthResponse{}, errors.New("user not found")
		} else {
			return dto.AuthResponse{}, err
		}
	}

	span.AddEvent("User found, verifying password")
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		span.AddEvent("Password verification failed")
		return dto.AuthResponse{}, err
	}

	claims := CustomClaims{
		Email:       loginRequest.Email,
		CurrentRole: loginRequest.CurrentRole,
	}

	authResponse, err := a.jwtService.GenerateToken(ctx, claims)

	if err != nil {
		span.RecordError(err)
		return authResponse, err
	}

	span.AddEvent("Storing refresh token in database")

	if user.RefreshToken.Token != "" {
		span.AddEvent("Invalidating existing refresh token")
		user.RefreshToken.IsValid = false
	}

	user.RefreshToken.Token = authResponse.RefreshToken
	user.RefreshToken.IsValid = true
	span.AddEvent("Updating user with new refresh token")

	user, err = a.userRepo.UpdateUser(ctx, user)

	if err != nil {
		span.RecordError(err)
		span.AddEvent("Failed to store refresh token")
		return dto.AuthResponse{}, err
	}

	span.AddEvent("Refresh token stored successfully")

	span.AddEvent("Mapping auth response")

	span.SetAttributes(attribute.String("user_id", user.Id.Hex()))

	return authResponse, nil
}

func (a *authService) Register(ctx context.Context, registerRequest dto.RegisterRequest) (dto.UserResponse, error) {

	ctx, span := utils.Tracer().Start(ctx, "AuthService.Register")
	defer span.End()

	existingUser, err := a.userRepo.FindUserByEmail(ctx, registerRequest.Email)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			span.AddEvent("User not found, proceeding with registration")
		} else {
			span.RecordError(err)
			return dto.UserResponse{}, err
		}
	}

	if existingUser.Email == registerRequest.Email {
		return dto.UserResponse{}, errors.New("user already exists")
	}

	span.AddEvent("Generating password hash")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)

	if err != nil {
		span.RecordError(err)
		return dto.UserResponse{}, err
	}

	now := time.Now()

	user := models.User{
		Name:        registerRequest.Name,
		Email:       registerRequest.Email,
		Password:    string(hashedPassword),
		PhoneNumber: registerRequest.PhoneNumber,
		Roles:       mappers.ToUserRoles(registerRequest.Roles),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	span.AddEvent("Creating user using repository")

	user, err = a.userRepo.CreateUser(ctx, user)

	if err != nil {
		span.RecordError(err)
		return dto.UserResponse{}, err
	}

	span.AddEvent("User created successfully")

	span.AddEvent("Mapping user to response")

	userResponse := mappers.ToUserResponse(user)

	span.AddEvent("User mapped to response successfully")

	return userResponse, nil
}
