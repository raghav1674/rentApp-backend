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

			rentRoutes := protectedRoutes.Group("/rents")
			landLordCheckMiddleWare := middlewares.RoleCheckMiddleware(string(models.LandLord))
			{
				rentRoutes.POST("",landLordCheckMiddleWare, rentController.CreateRent)
				rentRoutes.DELETE("/:rent_id",landLordCheckMiddleWare, rentController.CloseRent)
				rentRoutes.PUT("/:rent_id",landLordCheckMiddleWare, rentController.UpdateRent)
				rentRoutes.GET("", rentController.GetAllRents)
				rentRoutes.GET("/:rent_id", rentController.GetRentById)
			}
		}
	}
	return router
}
