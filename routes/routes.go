package routes

import (
	"sample-web/controllers"
	"sample-web/middlewares"
	"sample-web/models"
	"sample-web/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	userController controllers.UserController,
	authController controllers.AuthController,
	jwtService services.JWTService,
) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/v1")
	{

		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/login", authController.Login)
			authRoutes.POST("/register", authController.Register)
		}
		protectedRoutes := api.Group("/")
		protectedRoutes.Use(middlewares.JWTAuthMiddleware(jwtService))
		{
			userRoutes := protectedRoutes.Group("/users")
			{
				userRoutes.POST("/", userController.GetUserByEmail)
				userRoutes.PUT("/", userController.UpdateUser)
			}
			landlordRoutes := protectedRoutes.Group("/landlords")
			landlordRoutes.Use(middlewares.RoleCheckMiddleware(string(models.LandLord)))
			{
				landlordRoutes.POST("", userController.GetUserByEmail)
			}

			tenantRoutes := protectedRoutes.Group("/tenants")
			tenantRoutes.Use(middlewares.RoleCheckMiddleware(string(models.Tenant)))
			{
				tenantRoutes.POST("", userController.GetUserByEmail)
			}
		}
	}
	return router
}
