package util

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// FileUploadConfig holds configuration for file uploads
type FileUploadConfig struct {
	MaxFileSize       int64    // Maximum file size in bytes
	AllowedTypes      []string // Allowed MIME types
	AllowedExtensions []string // Allowed file extensions
	UploadDir         string   // Upload directory
}

// DefaultFileUploadConfig returns a secure default configuration
func DefaultFileUploadConfig() FileUploadConfig {
	return FileUploadConfig{
		MaxFileSize: 10 << 20, // 10MB
		AllowedTypes: []string{
			"image/jpeg",
			"image/png",
			"image/gif",
			"image/webp",
			"application/pdf",
		},
		AllowedExtensions: []string{
			".jpg", ".jpeg", ".png", ".gif", ".webp", ".pdf",
		},
		UploadDir: "uploads",
	}
}

// ValidateFile performs security checks on uploaded files
func ValidateFile(file *multipart.FileHeader, config FileUploadConfig) error {
	// Check file size
	if file.Size > config.MaxFileSize {
		return fiber.NewError(http.StatusBadRequest, "File too large")
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := false
	for _, allowedExt := range config.AllowedExtensions {
		if ext == allowedExt {
			allowed = true
			break
		}
	}
	if !allowed {
		return fiber.NewError(http.StatusBadRequest, "File type not allowed")
	}

	// Check MIME type
	fileContent, err := file.Open()
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Cannot read file")
	}
	defer fileContent.Close()

	// Read first 512 bytes to detect MIME type
	buffer := make([]byte, 512)
	n, err := fileContent.Read(buffer)
	if err != nil && err != io.EOF {
		return fiber.NewError(http.StatusInternalServerError, "Cannot read file content")
	}

	// Detect MIME type
	detectedType := http.DetectContentType(buffer[:n])

	// Additional security: Check for executable content in images
	if strings.HasPrefix(detectedType, "image/") {
		// Check for suspicious signatures in image files
		if isSuspiciousImage(buffer[:n]) {
			return fiber.NewError(http.StatusBadRequest, "Suspicious file content detected")
		}
	}

	// Check if detected MIME type is allowed
	allowed = false
	for _, allowedType := range config.AllowedTypes {
		if detectedType == allowedType {
			allowed = true
			break
		}
	}
	if !allowed {
		return fiber.NewError(http.StatusBadRequest, fmt.Sprintf("MIME type not allowed: %s", detectedType))
	}

	return nil
}

// isSuspiciousImage checks for potentially malicious content in image files
func isSuspiciousImage(buffer []byte) bool {
	// Check for script signatures that might be embedded
	suspiciousPatterns := []string{
		"<?php", "<?=", "<script", "javascript:", "vbscript:",
		"#!/bin/bash", "#!/bin/sh", "#!/usr/bin/env",
		"<%", "<%", "<%=", "<%@",
	}

	content := strings.ToLower(string(buffer))
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(content, strings.ToLower(pattern)) {
			return true
		}
	}

	// Check for null bytes that might indicate binary/script content
	if strings.Contains(content, "\x00") {
		return true
	}

	return false
}

// SaveUploadedFile saves an uploaded file with security checks
func SaveUploadedFile(file *multipart.FileHeader, subdir string, config *FileUploadConfig) (string, error) {
	// Use default config if nil
	if config == nil {
		defaultConfig := DefaultFileUploadConfig()
		config = &defaultConfig
	}

	// Validate file
	if err := ValidateFile(file, *config); err != nil {
		return "", err
	}

	// Create subdirectory if specified
	uploadPath := config.UploadDir
	if subdir != "" {
		// Sanitize subdirectory path
		subdir = filepath.Clean(subdir)
		// Prevent directory traversal
		if strings.Contains(subdir, "..") || strings.HasPrefix(subdir, "/") {
			return "", fiber.NewError(http.StatusBadRequest, "Invalid subdirectory")
		}
		uploadPath = filepath.Join(config.UploadDir, subdir)
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", fiber.NewError(http.StatusInternalServerError, "Cannot create upload directory")
	}

	// Generate unique filename to prevent overwrites
	ext := filepath.Ext(file.Filename)
	name := strings.TrimSuffix(file.Filename, ext)
	timestamp := time.Now().Format("20060102_150405_")
	filename := fmt.Sprintf("%s_%s%s", name, timestamp, ext)

	// Sanitize filename
	filename = strings.ReplaceAll(filename, "..", "")
	filename = strings.ReplaceAll(filename, "/", "")
	filename = strings.ReplaceAll(filename, "\\", "")

	filePath := filepath.Join(uploadPath, filename)

	// Save file directly (inline saveFile logic)
	src, err := file.Open()
	if err != nil {
		return "", fiber.NewError(http.StatusInternalServerError, "Cannot open uploaded file")
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fiber.NewError(http.StatusInternalServerError, "Cannot create destination file")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fiber.NewError(http.StatusInternalServerError, "Cannot save file")
	}

	// Return relative path for serving
	relativePath := filepath.Join(uploadPath, filename)
	return relativePath, nil
}

// GetFileInfo returns information about an uploaded file
func GetFileInfo(filePath string) (os.FileInfo, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, fiber.NewError(http.StatusNotFound, "File not found")
	}
	return info, nil
}

// DeleteFile removes an uploaded file
func DeleteFile(filePath string) error {
	// Validate path is within uploads directory to prevent directory traversal
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid file path")
	}

	uploadsDir, err := filepath.Abs("uploads")
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Cannot determine uploads directory")
	}

	// Ensure the file is within the uploads directory
	if !strings.HasPrefix(absPath, uploadsDir) {
		return fiber.NewError(http.StatusBadRequest, "File not in uploads directory")
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fiber.NewError(http.StatusNotFound, "File not found")
	}

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to delete file")
	}

	return nil
}
