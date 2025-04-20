package routes

import (
	"sample-web/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	userController controllers.UserController,
) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.POST("", userController.CreateUser)
			users.GET("", userController.GetUserByEmail)
			users.PUT("", userController.UpdateUser)
		}

		// rents := api.Group("/rents")
		// rents.GET("/", rentController.GetAllRents)
	}

	return router
}
