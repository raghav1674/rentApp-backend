package services

import (
	"context"
	"fmt"
	"sample-web/dto"
	"sample-web/utils"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	UserId      string `json:"user_id"`
	CurrentRole string `json:"current_role"`
}

type jwtCustomClaims struct {
	CustomClaims
	jwt.RegisteredClaims
}

type JWTService interface {
	GenerateAccessToken(ctx context.Context, customClaims CustomClaims) (string, error)
	GenerateRefreshToken(ctx context.Context) (string, error)
	GenerateToken(ctx context.Context, customClaims CustomClaims) (dto.AuthResponse, error)
	ValidateToken(ctx context.Context, token string) (*jwtCustomClaims, error)
}

type jwtService struct {
	secretKey                       string
	issuer                          string
	expirationInSeconds             int
	refreshTokenSecret              string
	refreshTokenExpirationInSeconds int
}

func NewJWTService(issuerName, secretKey, refreshTokenSecret string, refreshTokenExpirationInSeconds, expirationInSeconds int) JWTService {
	return &jwtService{
		issuer:                          issuerName,
		secretKey:                       secretKey,
		expirationInSeconds:             expirationInSeconds,
		refreshTokenSecret:              refreshTokenSecret,
		refreshTokenExpirationInSeconds: refreshTokenExpirationInSeconds,
	}
}

func (j *jwtService) GenerateAccessToken(ctx context.Context, customClaims CustomClaims) (string, error) {

	log := utils.GetLogger()

	spanCtx, span := log.Tracer().Start(ctx, "JWTService.GenerateToken")
	defer span.End()

	log.Info(spanCtx, fmt.Sprintf("Generating JWT token for user with user_id %s and current_role as %s", customClaims.UserId, customClaims.CurrentRole))

	claims := &jwtCustomClaims{
		CustomClaims: customClaims,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.expirationInSeconds) * time.Second)),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	log.Info(spanCtx, "Signing JWT token")

	return token.SignedString([]byte(j.secretKey))
}

func (j *jwtService) GenerateRefreshToken(ctx context.Context) (string, error) {

	log := utils.GetLogger()

	spanCtx, span := log.Tracer().Start(ctx, "JWTService.GenerateRefreshToken")
	defer span.End()

	log.Info(spanCtx, "Generating refresh token")

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.refreshTokenExpirationInSeconds) * time.Second)),
		Issuer:    j.issuer,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	log.Info(spanCtx, "Signing refresh token")

	return token.SignedString([]byte(j.refreshTokenSecret))
}

func (j *jwtService) ValidateToken(ctx context.Context, tokenString string) (*jwtCustomClaims, error) {

	log := utils.GetLogger()

	spanCtx, span := log.Tracer().Start(ctx, "JWTService.ValidateToken")
	defer span.End()

	log.Info(spanCtx, "Validating JWT token")

	var claims jwtCustomClaims

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		log.Info(spanCtx, "Parsing JWT token")
		return []byte(j.secretKey), nil
	})

	if err != nil {
		log.Error(spanCtx, err.Error())
		return nil, err
	}
	if claims, ok := token.Claims.(*jwtCustomClaims); ok && token.Valid {
		log.Info(spanCtx, "JWT token validated successfully")
		return claims, nil
	}
	log.Error(spanCtx, "JWT token validation failed")

	return nil, err
}

func (j *jwtService) GenerateToken(ctx context.Context, customClaims CustomClaims) (dto.AuthResponse, error) {

	log := utils.GetLogger()

	spanCtx, span := log.Tracer().Start(ctx, "JWTService.GenerateToken")
	defer span.End()

	log.Info(spanCtx, "Generating JWT Access token")

	accessToken, err := j.GenerateAccessToken(spanCtx, customClaims)
	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("Access Token generation failed with %s", err.Error()))
		return dto.AuthResponse{}, err
	}

	log.Info(spanCtx, "Access Token generated successfully")

	log.Info(spanCtx, "Generating JWT Refresh token")

	refreshToken, err := j.GenerateRefreshToken(spanCtx)

	if err != nil {
		log.Error(spanCtx, fmt.Sprintf("Refresh Token generation failed with %s", err.Error()))
		return dto.AuthResponse{}, err
	}
	log.Info(spanCtx, "Refresh Token generated successfully")
	return dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
