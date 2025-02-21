package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
    DBUrl      string
    JWTSecret  string
    Port       string
}

func LoadConfig() (*Config, error) {
    err := godotenv.Load()
    if err != nil {
        return nil, err
    }

    return &Config{
        DBUrl:     os.Getenv("DB_URL"),
        JWTSecret: os.Getenv("JWT_SECRET"),
        Port:      os.Getenv("PORT"),
    }, nil
}
