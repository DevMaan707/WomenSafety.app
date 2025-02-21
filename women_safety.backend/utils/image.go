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
	UploadDir   = "uploads/images"                       // Directory to store images
	MaxFileSize = 10 << 20                               // 10 MB
	BaseURL     = "http://localhost:3000/uploads/images" // Replace with your domain
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
	// Validate file size
	if file.Size > MaxFileSize {
		return "", &ImageUploadError{Message: "File size exceeds maximum limit"}
	}

	// Validate file type
	if !AllowedTypes[file.Header.Get("Content-Type")] {
		return "", &ImageUploadError{Message: "Invalid file type"}
	}

	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll(UploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s-%s%s",
		uuid.New().String(),
		time.Now().Format("20060102-150405"),
		ext,
	)

	// Create the file path
	filePath := filepath.Join(UploadDir, filename)

	// Open source
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %v", err)
	}
	defer src.Close()

	// Create destination
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file
	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file: %v", err)
	}

	// Return the URL
	return fmt.Sprintf("%s/%s", BaseURL, filename), nil
}

func DeleteImage(imageURL string) error {
	// Extract filename from URL
	filename := filepath.Base(imageURL)
	filePath := filepath.Join(UploadDir, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // File doesn't exist, nothing to delete
	}

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}

// Helper function to validate file type by extension
func isValidFileType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return true
	}
	return false
}
