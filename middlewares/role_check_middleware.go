package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleCheckMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println(c.Keys)
		currentRole, exists := c.Get("current_role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Current role not found"})
			return
		}

		currentRoleStr, ok := currentRole.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid role type"})
			return
		}

		if currentRoleStr != requiredRole {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized for this role"})
			return
		}

		c.Next()
	}
}
