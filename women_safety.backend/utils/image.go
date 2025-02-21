package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	UploadDir   = "uploads/images"
	MaxFileSize = 10 << 20
	BaseURL     = "http://localhost:3000/uploads/images"
)

var AllowedTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

type ImageUploadError struct {
	Message string
}

func (e *ImageUploadError) Error() string {
	return e.Message
}

func SaveImage(file *multipart.FileHeader) (string, error) {
	if file.Size > MaxFileSize {
		return "", &ImageUploadError{Message: "File size exceeds maximum limit"}
	}
	if !AllowedTypes[file.Header.Get("Content-Type")] {
		return "", &ImageUploadError{Message: "Invalid file type"}
	}
	if err := os.MkdirAll(UploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s-%s%s",
		uuid.New().String(),
		time.Now().Format("20060102-150405"),
		ext,
	)
	filePath := filepath.Join(UploadDir, filename)
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file: %v", err)
	}
	return fmt.Sprintf("%s/%s", BaseURL, filename), nil
}

func DeleteImage(imageURL string) error {
	filename := filepath.Base(imageURL)
	filePath := filepath.Join(UploadDir, filename)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	}
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}

func isValidFileType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return true
	}
	return false
}
