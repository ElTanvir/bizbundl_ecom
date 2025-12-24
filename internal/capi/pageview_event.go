package capi

import (
	"crypto/sha256"
	"portfolioed/internal/store"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func SendPageViewEvent(c *fiber.Ctx) {
	pixelID := store.GetFBPixel()
	accessToken := store.GetFBPixelToken()
	if pixelID == "" || accessToken == "" {
		return
	}
	fbpCookie := c.Cookies("_fbp")
	fbcCookie := c.Cookies("_fbc")
	ipAddress := c.IP()
	userAgent := c.Get("User-Agent")
	fullURL := c.Protocol() + "://" + c.Hostname() + c.OriginalURL()
	eventData := CAPIEventData{
		FBP:       fbpCookie,
		FBC:       fbcCookie,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		FullURL:   fullURL,
	}
	go sendPageViewEventAsync(pixelID, accessToken, eventData)
}
func sendPageViewEventAsync(pixelID, accessToken string, eventData CAPIEventData) {
	userData := UserData{
		ClientIPAddress: eventData.IPAddress,
		ClientUserAgent: eventData.UserAgent,
		FBC:             eventData.FBC,
		FBP:             eventData.FBP,
	}
	event := CAPIEvent{
		EventName:      "PageView",
		EventTime:      time.Now().Unix(),
		ActionSource:   "website",
		EventID:        uuid.New().String(),
		EventSourceURL: eventData.FullURL,
		UserData:       userData,
	}

	payload := map[string]any{
		"data": []CAPIEvent{event},
		// "test_event_code": "TEST53511",
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Error().Msgf("Error marshaling CAPI payload: %v", err)
		return
	}
	sendCAPIRequest(pixelID, accessToken, jsonPayload)
}

func HashValue(value string) string {
	hasher := sha256.New()
	hasher.Write([]byte(value))
	return hex.EncodeToString(hasher.Sum(nil))
}
