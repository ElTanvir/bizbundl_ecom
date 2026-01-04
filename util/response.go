package util

import (
	"github.com/gofiber/fiber/v2"
)

// APIResponse is the standard structure for all API responses
type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// JSON sends a success response
func JSON(c *fiber.Ctx, status int, data any, message string) error {
	return c.Status(status).JSON(APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// APIError sends an error response
func APIError(c *fiber.Ctx, status int, err error) error {
	msg := "Unknown error"
	if err != nil {
		msg = err.Error()
	}
	return c.Status(status).JSON(APIResponse{
		Success: false,
		Error:   msg,
	})
}
