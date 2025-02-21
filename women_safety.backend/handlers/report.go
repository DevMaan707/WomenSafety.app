package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"DevMaan707/Saathee/database"
	"DevMaan707/Saathee/utils"
)

func PostReport(c *fiber.Ctx) error {
	// Parse form data
	report := new(database.Report)
	if err := c.BodyParser(report); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Handle file upload
	file, err := c.FormFile("image")
	if err == nil {
		// File was provided, save it
		imageURL, err := utils.SaveImage(file)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		report.ImageURL = imageURL
	}

	report.ID = utils.GenerateUUID()
	report.ReportedBy = c.Locals("user_id").(string)
	report.Status = "PENDING"
	report.CreatedAt = time.Now()

	// Insert into database
	_, err = database.DB.Exec(`
        INSERT INTO reports (id, image_url, latitude, longitude, description, reported_by, status, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		report.ID, report.ImageURL, report.Latitude, report.Longitude, report.Description,
		report.ReportedBy, report.Status, report.CreatedAt)
	if err != nil {
		// If there was an error saving to database, delete the uploaded image
		if report.ImageURL != "" {
			_ = utils.DeleteImage(report.ImageURL)
		}
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create report",
		})
	}

	return c.Status(201).JSON(report)
}

// Optional: Add function to serve static files
func SetupStaticFiles(app *fiber.App) {
	app.Static("/uploads", "./uploads")
}

func GetReports(c *fiber.Ctx) error {
	if c.Locals("role") != string(database.RoleAuthority) {
		return c.Status(403).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var reports []database.Report
	err := database.DB.Select(&reports, "SELECT * FROM reports ORDER BY created_at DESC")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch reports",
		})
	}

	return c.JSON(reports)
}
