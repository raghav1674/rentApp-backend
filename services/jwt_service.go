package services

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
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
	GenerateToken(ctx *gin.Context, claims CustomClaims) (string, error)
	ValidateToken(ctx *gin.Context, token string) (*jwtCustomClaims, error)
}

type jwtService struct {
	secretKey           string
	issuer              string
	expirationInSeconds int
}

func NewJWTService(issuerName, secretKey string, expirationInSeconds int) JWTService {
	return &jwtService{
		issuer:              issuerName,
		secretKey:           secretKey,
		expirationInSeconds: expirationInSeconds,
	}
}

func (j *jwtService) GenerateToken(ctx *gin.Context, customClaims CustomClaims) (string, error) {
	claims := &jwtCustomClaims{
		CustomClaims: customClaims,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.expirationInSeconds) * time.Second)),
			Issuer:    j.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

func (j *jwtService) ValidateToken(ctx *gin.Context, tokenString string) (*jwtCustomClaims, error) {

	var claims jwtCustomClaims

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*jwtCustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
