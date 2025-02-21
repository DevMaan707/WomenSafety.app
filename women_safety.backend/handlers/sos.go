package handlers

import (
	"fmt"
	"time"

	"github.com/MelloB1989/karma/apis/aws/sns"
	d "github.com/MelloB1989/karma/database"
	"github.com/MelloB1989/karma/orm"
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
	go notifyAuthorities()
	fmt.Println("SOS Alert Triggered")

	return c.Status(201).JSON(sos)
}

func notifyAuthorities() {
	db, err := d.PostgresConn()
	if err != nil {
		fmt.Println("Failed to connect to database")
		return
	}

	query := "SELECT sns_arn FROM users WHERE sns_arn IS NOT NULL"

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Failed to fetch authorities")
		return
	}

	for rows.Next() {
		var snsARN string
		err = rows.Scan(&snsARN)
		if err != nil {
			fmt.Println("Failed to scan authorities")
			return
		}

		// Send notification
		go notify(snsARN)
	}
}

func notify(snsarn string) {
	notification := `{
  "GCM": "{\"notification\":{\"title\":\"%s\",\"body\":\"%s\"}}"
  }`

	karmasns := sns.NewFCM()
	fmt.Println(karmasns.ListAllEndpointARNs())
	karmasns.PublishGCMMessageToAllEndpoints(fmt.Sprintf(notification, "SOS", "Laude ka SOS, it was just a prank"))
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

func RegisterSNS(c *fiber.Ctx) error {

	type requestBody struct {
		SnsArn string `json:"sns_arn"`
	}

	req := new(requestBody)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	uid := c.Locals("user_id").(string)

	userORM := orm.Load(&database.User{})

	u, err := userORM.GetByFieldEquals("Id", uid)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	user, ok := u.([]*database.User)
	if !ok {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch user",
		})
	}

	user[0].SnsArn = req.SnsArn

	err = userORM.Update(user[0], user[0].Id)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "SNS ARN registered successfully",
	})

}
