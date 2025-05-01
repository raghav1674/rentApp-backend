package middlewares

import (
	"net/http"
	"sample-web/utils"

	"github.com/gin-gonic/gin"
)

func RoleCheckMiddleware(requiredRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		_, span := utils.Tracer().Start(ctx.Request.Context(), "middlewares.RoleCheckMiddleware")
		defer span.End()

		currentRole, exists := ctx.Get("current_role")

		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Current role not found"})
			return
		}

		currentRoleStr, ok := currentRole.(string)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid role type"})
			return
		}

		if currentRoleStr != requiredRole {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized for this role"})
			return
		}

		ctx.Next()
	}
}
