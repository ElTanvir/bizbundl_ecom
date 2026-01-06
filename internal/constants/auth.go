package constants

import "time"

const (
	UserSessionDuration  = 7 * 24 * time.Hour
	GuestSessionDuration = 2 * 365 * 24 * time.Hour
	RefreshThreshold     = 30 * time.Minute
)
