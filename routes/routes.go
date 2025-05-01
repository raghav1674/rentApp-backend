package routes

import (
	"sample-web/controllers"
	"sample-web/middlewares"
	"sample-web/models"
	"sample-web/services"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func SetupRouter(
	userController controllers.UserController,
	authController controllers.AuthController,
	jwtService services.JWTService,
) *gin.Engine {
	router := gin.Default()
	router.Use(otelgin.Middleware("sample-web"))
	router.Use(middlewares.ErrorHandler())

	api := router.Group("/api/v1")
	{

		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", authController.Register)
			authRoutes.POST("/login", authController.Login)
		}
		protectedRoutes := api.Group("/")
		protectedRoutes.Use(middlewares.JWTAuthMiddleware(jwtService))
		{
			userRoutes := protectedRoutes.Group("/users")
			{
				userRoutes.POST("", userController.GetUserByEmail)
				userRoutes.PUT("", userController.UpdateUser)
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
