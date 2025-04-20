package main

import (
	"sample-web/clients"
	"sample-web/configs"
	"sample-web/controllers"
	"sample-web/repositories"
	"sample-web/routes"
	"sample-web/services"
)

const (
	configPath = "config.json"
)

func main() {

	appConfigs, err := configs.LoadConfig(configPath)

	if err != nil {
		panic(err)
	}

	mongoClient, err := clients.NewMongoClient(appConfigs.GetMongoConfig())

	if err != nil {
		panic(err)
	}

	// Wire up dependencies
	userRepo := repositories.NewUserRepository(mongoClient.Database)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	// Set up router with all routes
	r := routes.SetupRouter(userController)

	// Start the server
	r.Run(":8080")
}
