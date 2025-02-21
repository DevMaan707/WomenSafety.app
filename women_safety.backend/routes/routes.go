package routes

import (
	"github.com/gofiber/fiber/v2"

	"DevMaan707/Saathee/database"
	"DevMaan707/Saathee/handlers"
	"DevMaan707/Saathee/middleware"
)

func SetupRoutes(app *fiber.App) {
	// Public routes

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(database.ResponseHTTP{
			Success: true,
			Data:    nil,
			Message: "OK",
		})
	})

	app.Post("/register", handlers.Register)
	app.Post("/login", handlers.Login)

	// Protected routes
	api := app.Group("/api", middleware.AuthMiddleware)

	// Location routes
	api.Get("/locations", handlers.GetLocations)
	api.Post("/locations", handlers.PostLocation)

	// Report routes
	api.Post("/reports", handlers.PostReport)
	api.Get("/reports", handlers.GetReports)

	// SOS routes
	api.Post("/sos", handlers.TriggerSOS)
	api.Get("/sos", handlers.GetSOSAlerts)
}
