package main

import (
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
	redisConfig := appConfigs.GetRedisConfig()
	jwtConfig := appConfigs.GetJWTConfig()
	tracingConfig := appConfigs.GetTracingConfig()
	twilioConfig := appConfigs.GetTwilioConfig()

	utils.InitLogger(tracingConfig)

	// Initialize MongoDB client
	mongoClient, err := clients.NewMongoClient(mongoConfig)
	if err != nil {
		panic(err)
	}

	redisClient, err := clients.NewRedisClient(redisConfig)
	if err != nil {
		panic(err)
	}

	// Initialize JWT service
	jwtService := services.NewJWTService(jwtConfig.IssuerName,
		jwtConfig.SecretKey,
		jwtConfig.RefreshTokenSecret,
		jwtConfig.ExpirationInSeconds,
		jwtConfig.RefreshTokenExpirationInSeconds)

	// Initialize OTP service
	otpService := services.NewTwilioOTPService(twilioConfig, redisClient)
	// otpService := services.NewDummyOTPService(redisClient)

	// Initialize the user repository, service, and controller
	userRepo := repositories.NewUserRepository(mongoClient.Database)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	// Initialize the auth service and controller
	authService := services.NewAuthService(userRepo, jwtService)
	authController := controllers.NewAuthController(authService, otpService)

	// Initialize rent repository, service, and controller
	rentRepo := repositories.NewRentRepository(mongoClient.Database)
	rentService := services.NewRentService(rentRepo, userRepo)
	rentController := controllers.NewRentController(rentService)

	// initialize rent record	repository, service, and controller
	rentRecordRepo := repositories.NewRentRecordRepository(mongoClient.Database)
	rentRecordService := services.NewRentRecordService(rentRecordRepo, rentRepo, userRepo)
	rentRecordController := controllers.NewRentRecordController(rentRecordService)

	// Initialize the health controller
	healthController := controllers.NewHealthController()

	// Set up router with all routes
	r := routes.SetupRouter(healthController,userController, authController, rentController, rentRecordController, jwtService)
	// Start the server
	r.Run(":8080")
}
