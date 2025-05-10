package middlewares

import (
	"fmt"
	"net/http"
	"sample-web/services"
	"sample-web/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware(jwtService services.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		log := utils.GetLogger()

		spanCtx, span := log.Tracer().Start(ctx.Request.Context(), "middlewares.JWTAuthMiddleware")
		defer span.End()

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

		customClaims, err := jwtService.ValidateToken(spanCtx, tokenString[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		ctx.Set("user_id", customClaims.UserId)
		ctx.Set("current_role", customClaims.CurrentRole)
		fmt.Print("User ID: ", customClaims.UserId)
		ctx.Next()
	}
}
