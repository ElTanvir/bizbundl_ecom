package capi

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)


func sendCAPIRequest(pixelID, accessToken string, jsonPayload []byte) {
	agent := fiber.Post(fmt.Sprintf("https://graph.facebook.com/v19.0/%s/events?access_token=%s", pixelID, accessToken))
	agent.Body(jsonPayload)
	agent.Set("Content-Type", "application/json")
	agent.Timeout(10 * time.Second)
	statusCode, _, errs := agent.Bytes()
	if len(errs) > 0 {
		log.Error().Msgf("Error sending CAPI request with Fiber: %v", errs[0])
		return
	}
	if statusCode != fiber.StatusOK {
		log.Error().Msgf("CAPI request failed with status: %d", statusCode)
	}
}
