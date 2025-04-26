package middlewares

import (
	"net/http"
	"sample-web/services"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(jwtService services.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		tokenString := strings.Split(authHeader, "Bearer ")
		if len(tokenString) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			return
		}

		customClaims, err := jwtService.ValidateToken(ctx, tokenString[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		ctx.Set("email", customClaims.Email)
		ctx.Set("current_role", customClaims.CurrentRole)

		ctx.Next()
	}
}
