package main

import (
	"context"
	"os"
	"sample-web/clients"
	"sample-web/configs"
	"sample-web/controllers"
	"sample-web/repositories"
	"sample-web/routes"
	"sample-web/services"
	"sample-web/utils"
)

const (
	defaultConfigPath = "config.json"
)

func main() {

	// Load configuration
	configPath := os.Getenv("CONFIG_FILE_PATH")
	if configPath == "" {
		configPath = defaultConfigPath
	}
	appConfigs, err := configs.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}
	mongoConfig := appConfigs.GetMongoConfig()
	jwtConfig := appConfigs.GetJWTConfig()
	tracingConfig := appConfigs.GetTracingConfig()

	shutdown := utils.InitTracer(tracingConfig)
	defer shutdown(context.Background())

	// Initialize MongoDB client
	mongoClient, err := clients.NewMongoClient(mongoConfig)
	if err != nil {
		panic(err)
	}

	// Initialize JWT service
	jwtService := services.NewJWTService(jwtConfig.IssuerName, jwtConfig.SecretKey, jwtConfig.ExpirationInSeconds)

	// Initialize the user repository, service, and controller
	userRepo := repositories.NewUserRepository(mongoClient.Database)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	// Initialize the auth service and controller
	authService := services.NewAuthService(userRepo, jwtService)
	authController := controllers.NewAuthController(authService)

	// Set up router with all routes
	r := routes.SetupRouter(userController, authController, jwtService)

	// Start the server
	r.Run(":8080")
}
