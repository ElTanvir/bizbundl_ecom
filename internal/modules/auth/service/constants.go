package service

import "time"

const (
	// UserSessionDuration is the core validity of a logged-in user token.
	// User asked for "Short". Let's do 2 Hours.
	// With auto-renewal, this means 2 hours of *inactivity* logs you out.
	UserSessionDuration = 2 * time.Hour

	// GuestSessionDuration is 2 years to persist carts.
	GuestSessionDuration = 2 * 365 * 24 * time.Hour

	// RefreshThreshold is the elapsed time after which we auto-renew the token.
	// E.g., if session is 2h, and user has been active for 30 mins, we renew back to full 2h.
	RefreshThreshold = 30 * time.Minute
)
