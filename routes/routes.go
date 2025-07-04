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
	healthController controllers.HealthController,
	userController controllers.UserController,
	authController controllers.AuthController,
	rentController controllers.RentController,
	rentRecordController controllers.RentRecordController,
	jwtService services.JWTService,
) *gin.Engine {
	router := gin.Default()
	router.Use(otelgin.Middleware("sample-web"))
	router.Use(middlewares.ErrorHandler())



	api := router.Group("/api/v1")
	{
		api.GET("/health", healthController.GetHealth)
		
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
				userRoutes.GET("/me", userController.GetCurrentUser)
				userRoutes.POST("", userController.GetUserByPhoneNumber)
				userRoutes.PUT("", userController.UpdateUser)
			}

			landLordCheckMiddleWare := middlewares.RoleCheckMiddleware(string(models.LandLord))
			tenantCheckMiddleWare := middlewares.RoleCheckMiddleware(string(models.Tenant))

			rentRoutes := protectedRoutes.Group("/rents")
			{
				rentRoutes.POST("", landLordCheckMiddleWare, rentController.CreateRent)
				rentRoutes.DELETE("/:rent_id", landLordCheckMiddleWare, rentController.CloseRent)
				rentRoutes.PUT("/:rent_id", landLordCheckMiddleWare, rentController.UpdateRent)
				rentRoutes.GET("", rentController.GetAllRents)
				rentRoutes.GET("/:rent_id", rentController.GetRentById)
			}
			rentRecordRoutes := protectedRoutes.Group("/rents/:rent_id/records")
			{
				rentRecordRoutes.POST("", tenantCheckMiddleWare, rentRecordController.CreateRentRecord)
				rentRecordRoutes.GET("", rentRecordController.GetAllRentRecords)
				rentRecordRoutes.GET("/:record_id", rentRecordController.GetRentRecordById)
				rentRecordRoutes.POST("/:record_id/approve", landLordCheckMiddleWare, rentRecordController.ApproveRentRecord)
				rentRecordRoutes.POST("/:record_id/reject", landLordCheckMiddleWare, rentRecordController.RejectRentRecord)
			}
		}
	}
	return router
}
