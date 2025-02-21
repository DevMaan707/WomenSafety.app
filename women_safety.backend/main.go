package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"DevMaan707/Saathee/config"
	"DevMaan707/Saathee/database"
	"DevMaan707/Saathee/routes"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}
	err = database.InitDB(cfg.DBUrl)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024,
	})
	app.Static("/uploads", "./uploads")
	routes.SetupRoutes(app)
	log.Fatal(app.Listen(":" + cfg.Port))
}
