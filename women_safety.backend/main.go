package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"DevMaan707/Saathee/config"
	"DevMaan707/Saathee/database"
	"DevMaan707/Saathee/routes"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize database
	err = database.InitDB(cfg.DBUrl)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024, // 10MB limit for file uploads
	})

	// Setup static file serving
	app.Static("/uploads", "./uploads")

	// Setup routes
	routes.SetupRoutes(app)

	// Start server
	log.Fatal(app.Listen(":" + cfg.Port))
}
