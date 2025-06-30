package main

import (
	_ "github.com/hadisjane/confessly/docs" // Import the generated docs package
	"github.com/hadisjane/confessly/internal/configs"
	"github.com/hadisjane/confessly/internal/controller"
	"github.com/hadisjane/confessly/internal/db"
	"github.com/hadisjane/confessly/logger"
	"log"
)

// @title Confessly API
// @version 1.0
// @description API Server for Confessly Application
// @securityDefinitions.apikey ApiKeyAuth
// @host localhost:8081
// @BasePath /
// @in header
// @name Authorization

func main() {
	// Load configurations
	if err := configs.ReadSettings(); err != nil {
		log.Fatalf("Failed to load configurations: %v", err)
	}

	// Initialize logger
	if err := logger.Init(); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Initialize database connection
	if err := db.ConnDB(); err != nil {
		logger.Error.Fatalf("Error connecting to database: %v", err)
	}
	logger.Info.Println("Database connection established successfully")

	// Run schema initialization
	if err := db.InitMigrations(); err != nil {
		logger.Error.Fatalf("Error initializing database schema: %v", err)
	}
	logger.Info.Println("Database schema initialized successfully")

	// Seed database
	if err := db.SeedDB(); err != nil {
		logger.Error.Fatalf("Error seeding database: %v", err)
	}
	logger.Info.Println("Database seeded successfully")

	// Start the server
	if err := controller.RunServer(); err != nil {
		logger.Error.Fatalf("Error running server: %v", err)
	}
}
