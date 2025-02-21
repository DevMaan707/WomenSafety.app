package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"DevMaan707/Saathee/database"
	"DevMaan707/Saathee/utils"
)

func GetLocations(c *fiber.Ctx) error {
	var locations []database.RiskLocation
	err := database.DB.Select(&locations, "SELECT * FROM risk_locations")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch locations",
		})
	}

	return c.JSON(locations)
}

func PostLocation(c *fiber.Ctx) error {
	location := new(database.RiskLocation)
	if err := c.BodyParser(location); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	location.ID = utils.GenerateUUID()
	location.CreatedBy = c.Locals("user_id").(string)
	location.CreatedAt = time.Now()

	_, err := database.DB.Exec(
		`
        INSERT INTO risk_locations (id, latitude, longitude, risk_level, created_by, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)`,
		location.ID,
		location.Latitude,
		location.Longitude,
		location.RiskLevel,
		location.CreatedBy,
		location.CreatedAt,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create location",
		})
	}

	return c.Status(201).JSON(location)
}
