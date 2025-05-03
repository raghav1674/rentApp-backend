package services

import (
	"context"
	"sample-web/dto"
	"sample-web/utils"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type CustomClaims struct {
	Email       string `json:"email"`
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

	_, span := utils.Tracer().Start(ctx, "JWTService.GenerateToken")
	defer span.End()

	span.SetAttributes(attribute.String("current_role", customClaims.CurrentRole))

	span.AddEvent("Generating JWT token", trace.WithAttributes(
		attribute.String("email", customClaims.Email),
		attribute.String("current_role", customClaims.CurrentRole),
	))

	claims := &jwtCustomClaims{
		CustomClaims: customClaims,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.expirationInSeconds) * time.Second)),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	span.AddEvent("Signing JWT token")

	return token.SignedString([]byte(j.secretKey))
}

func (j *jwtService) GenerateRefreshToken(ctx context.Context) (string, error) {

	_, span := utils.Tracer().Start(ctx, "JWTService.GenerateRefreshToken")
	defer span.End()

	span.AddEvent("Generating refresh token")

	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.refreshTokenExpirationInSeconds) * time.Second)),
		Issuer:    j.issuer,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	span.AddEvent("Signing refresh token")

	return token.SignedString([]byte(j.refreshTokenSecret))
}

func (j *jwtService) ValidateToken(ctx context.Context, tokenString string) (*jwtCustomClaims, error) {

	_, span := utils.Tracer().Start(ctx, "JWTService.ValidateToken")
	defer span.End()

	span.AddEvent("Validating JWT token")

	var claims jwtCustomClaims

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		span.AddEvent("Parsing JWT token")
		return []byte(j.secretKey), nil
	})

	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	if claims, ok := token.Claims.(*jwtCustomClaims); ok && token.Valid {
		span.AddEvent("JWT token validated successfully")
		return claims, nil
	}
	span.RecordError(err)
	span.AddEvent("JWT token validation failed")
	return nil, err
}

func (j *jwtService) GenerateToken(ctx context.Context, customClaims CustomClaims) (dto.AuthResponse, error) {

	spanCtx, span := utils.Tracer().Start(ctx, "JWTService.GenerateToken")
	defer span.End()

	span.AddEvent("Generating JWT Access token")

	accessToken, err := j.GenerateAccessToken(spanCtx, customClaims)
	if err != nil {
		span.RecordError(err)
		span.AddEvent("Access Token generation failed")
		return dto.AuthResponse{}, err
	}

	span.AddEvent("Access Token generated successfully")

	span.AddEvent("Generating JWT Refresh token")

	refreshToken, err := j.GenerateRefreshToken(spanCtx)

	if err != nil {
		span.RecordError(err)
		span.AddEvent("Refresh Token generation failed")
		return dto.AuthResponse{}, err
	}
	span.AddEvent("Refresh Token generated successfully")
	return dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
