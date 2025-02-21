package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"DevMaan707/Saathee/config"
	"DevMaan707/Saathee/database"
	"DevMaan707/Saathee/utils"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Language string `json:"language"`
	Gender   string `json:"gender"`
	Aadhaar  string `json:"aadhaar"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func Register(c *fiber.Ctx) error {
	req := new(RegisterRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	user := database.User{
		ID:        utils.GenerateUUID(),
		Name:      req.Name,
		Age:       req.Age,
		Language:  req.Language,
		Gender:    req.Gender,
		Aadhaar:   req.Aadhaar,
		Password:  string(hashedPassword),
		Role:      database.Role(req.Role),
		CreatedAt: time.Now(),
	}

	_, err = database.DB.Exec(
		`
        INSERT INTO users (id, name, age, language, gender, aadhaar, password, role, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		user.ID,
		user.Name,
		user.Age,
		user.Language,
		user.Gender,
		user.Aadhaar,
		user.Password,
		user.Role,
		user.CreatedAt,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(201).JSON(database.ResponseHTTP{
		Success: true,
		Data:    user.ID,
		Message: "User created successfully",
	})
}

func Login(c *fiber.Ctx) error {
	var credentials struct {
		Aadhaar  string `json:"aadhaar"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&credentials); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	var user database.User
	err := database.DB.Get(&user, "SELECT * FROM users WHERE aadhaar = $1", credentials.Aadhaar)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	conf, err := config.LoadConfig()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to load config",
		})
	}

	tokenString, err := token.SignedString([]byte(conf.JWTSecret))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(fiber.Map{
		"token": tokenString,
	})
}
