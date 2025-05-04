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
	rentController controllers.RentController,
	jwtService services.JWTService,
) *gin.Engine {
	router := gin.Default()
	router.Use(otelgin.Middleware("sample-web"))
	router.Use(middlewares.ErrorHandler())

	api := router.Group("/api/v1")
	{

		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/otp/generate", authController.GenerateOTP)
			authRoutes.POST("/otp/verify", authController.VerifyOTP)
			authRoutes.POST("/register", authController.Register)
		}
		protectedRoutes := api.Group("/")
		protectedRoutes.Use(middlewares.JWTAuthMiddleware(jwtService))
		{
			userRoutes := protectedRoutes.Group("/users")
			{
				userRoutes.GET("/me",userController.GetCurrentUser)
				userRoutes.POST("", userController.GetUserByPhoneNumber)
				userRoutes.PUT("", userController.UpdateUser)
			}
			landlordRoutes := protectedRoutes.Group("/landlords")
			landlordRoutes.Use(middlewares.RoleCheckMiddleware(string(models.LandLord)))
			{
				landlordRoutes.POST("", userController.GetUserByPhoneNumber)
				landlordRentRoutes := landlordRoutes.Group("/rents")
				{
					landlordRentRoutes.POST("", rentController.CreateRent)
					landlordRentRoutes.GET("", rentController.GetAllRents)
					landlordRentRoutes.GET("/:rent_id", rentController.GetRentById)
					landlordRentRoutes.PUT("/:rent_id", rentController.UpdateRent)
					landlordRentRoutes.DELETE("/:rent_id", rentController.CloseRent)
				}
			}

			tenantRoutes := protectedRoutes.Group("/tenants")
			tenantRoutes.Use(middlewares.RoleCheckMiddleware(string(models.Tenant)))
			{
				tenantRentRoutes := tenantRoutes.Group("/rents")
				{
					tenantRentRoutes.GET("", rentController.GetAllRents)
					tenantRentRoutes.GET("/:rent_id", rentController.GetRentById)
				}
			}
		}
	}
	return router
}
