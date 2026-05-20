package main

import (
	"log"
	"os"
	"studytrack/config"
	"studytrack/models"
	"studytrack/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Koneksi database
	config.ConnectDatabase()

	// Auto migrate
	config.DB.AutoMigrate(
		&models.User{},
		&models.OTP{},
		&models.Pomodoro{},
		&models.DailyTarget{},
    	&models.TargetMatkul{},

	)

	// Setup Gin
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r)

	// Jalankan server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}

