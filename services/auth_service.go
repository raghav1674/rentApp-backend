package services

import (
	"sample-web/dto"
	"sample-web/mappers"
	"sample-web/models"
	"sample-web/repositories"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx *gin.Context, loginRequest dto.LoginRequest) (dto.AuthResponse, error)
	Register(ctx *gin.Context, registerRequest dto.RegisterRequest) (dto.UserResponse, error)
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

func (a *authService) Login(ctx *gin.Context, loginRequest dto.LoginRequest) (dto.AuthResponse, error) {
	user, err := a.userRepo.FindUserByEmail(ctx, loginRequest.Email)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		return dto.AuthResponse{}, err
	}

	claims := CustomClaims{
		Email:       loginRequest.Email,
		CurrentRole: loginRequest.CurrentRole,
	}

	accesToken, err := a.jwtService.GenerateToken(ctx, claims)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	authResponse := dto.AuthResponse{
		AccessToken: accesToken,
	}
	return authResponse, nil
}

func (a *authService) Register(ctx *gin.Context, registerRequest dto.RegisterRequest) (dto.UserResponse, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.UserResponse{}, err
	}

	now := time.Now()

	user := models.User{
		Email:       registerRequest.Email,
		Password:    string(hashedPassword),
		PhoneNumber: registerRequest.PhoneNumber,
		Roles:       mappers.ToUserRoles(registerRequest.Roles),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	user, err = a.userRepo.CreateUser(ctx, user)
	if err != nil {
		return dto.UserResponse{}, err
	}
	userResponse := mappers.ToUserResponse(user)
	return userResponse, nil
}
