package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	"DevMaan707/Saathee/database"
	"DevMaan707/Saathee/utils"
)

func TriggerSOS(c *fiber.Ctx) error {
	sos := new(database.SOS)
	if err := c.BodyParser(sos); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	sos.ID = utils.GenerateUUID()
	sos.UserID = c.Locals("user_id").(string)
	sos.Active = true
	sos.CreatedAt = time.Now()

	_, err := database.DB.Exec(`
        INSERT INTO sos (id, user_id, latitude, longitude, active, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)`,
		sos.ID, sos.UserID, sos.Latitude, sos.Longitude, sos.Active, sos.CreatedAt)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create SOS alert",
		})
	}

	// Notify authorities and nearby users
	// go notifyAuthorities(sos)
	// go notifyNearbyUsers(sos)
	fmt.Println("SOS Alert Triggered")

	return c.Status(201).JSON(sos)
}

func GetSOSAlerts(c *fiber.Ctx) error {
	if c.Locals("role") != string(database.RoleAuthority) {
		return c.Status(403).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var alerts []database.SOS
	err := database.DB.Select(
		&alerts,
		"SELECT * FROM sos WHERE active = true ORDER BY created_at DESC",
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch SOS alerts",
		})
	}

	return c.JSON(alerts)
}
