package main

import (
	"log"

	"otp-service/config"
	"otp-service/internal/api"
	"otp-service/internal/repository"
	"otp-service/internal/service"
	"otp-service/pkg/database"
	"otp-service/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func main() {

	conf := config.Load()

	// init databse
	db, err := database.NewConnection(conf.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database :%v", err)
	}

	// Migrate database
	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed migrations :%v", err)
	}

	// init repositories
	userRepo := repository.NewUserRepository(db)
	otpRepo := repository.NewOTPRepository()

	// init services
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, otpRepo, conf.JWTSecret)

	// init handlers
	userHandler := api.NewUserHandler(userService)
	authHandler := api.NewAuthHandler(authService)

	// setting up gin router
	router := gin.Default()

	// add middle ware
	router.Use(middleware.CORS())
	router.Use(middleware.RequestLogger())

	// setup routes

	api.SetupRoutes(router, authHandler, userHandler, conf.JWTSecret)

	// start server
	log.Printf("*Server start on port %s", conf.Port)
	if err := router.Run(":" + conf.Port); err != nil {
		log.Fatalf("Failed to start server : %v", err)
	}

}
