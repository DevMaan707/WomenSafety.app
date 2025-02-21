package utils

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func GenerateUUID() string {
	id, _ := gonanoid.Generate("qwertyuiopasdfghjklzxcvbnm1234567890_-", 10)
	return id
}
