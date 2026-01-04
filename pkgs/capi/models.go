package capi
type CAPIEventData struct {
	FBP       string
	FBC       string
	IPAddress string
	UserAgent string
	FullURL   string
}
type CAPIEvent struct {
	EventName      string   `json:"event_name"`
	EventTime      int64    `json:"event_time"`
	EventSourceURL string   `json:"event_source_url"`
	EventID        string   `json:"event_id"`
	UserData       UserData `json:"user_data"`
	ActionSource   string   `json:"action_source"`
}

// UserData contains all the user-specific identifiers.
type UserData struct {
	ClientIPAddress string   `json:"client_ip_address,omitempty"`
	ClientUserAgent string   `json:"client_user_agent,omitempty"`
	FBC             string   `json:"fbc,omitempty"`
	FBP             string   `json:"fbp,omitempty"`
	Email           []string `json:"em,omitempty"`
	Phone           []string `json:"ph,omitempty"`
	ExternalID      []string `json:"external_id,omitempty"`
}
