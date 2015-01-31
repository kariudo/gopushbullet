package gopushbullet

//ErrorResponse (any non-200 error code) contain information on the kind of error that happened.
type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Cat     string `json:"cat"`
	} `json:"error"`
}

//PushMessage describes a message to be sent via PushBullet. Only one of the first 4 properties may be specified with a message being sent.
type PushMessage struct {
	DeviceID   string `json:"device_iden"`
	Email      string `json:"email"`
	ChannelTag string `json:"channel_tag"`
	ClientID   string `json:"client_iden"`

	// Type indicates the type of push message being sent
	Type    string   `json:"type"`    // note, link, address, checklist, or file
	Title   string   `json:"title"`   // Title is used withn note, link, and checklist types
	Body    string   `json:"body"`    // Body is used with link, note, or file types
	URL     string   `json:"url"`     // URL is used for link types
	Address string   `json:"address"` // Address is used with address type
	Items   []string `json:"items"`   // Items are used with checklist types
	// The following are used with file types
	FileName       string `json:"file_name"`
	FileType       string `json:"file_type"` // MIME type of the file
	FileURL        string `json:"file_url"`
	SourceDeviceID string `json:"source_device_iden"`
}

//Device describes a registered device (phone, stream).
type Device struct {
	ID           string  `json:"iden"`
	PushToken    string  `json:"push_token"`
	AppVersion   int     `json:"app_version"`
	Fingerprint  string  `json:"fingerprint"`
	Active       bool    `json:"active"`
	Nickname     string  `json:"nickname"`
	Manufacturer string  `json:"manufacturer"`
	Type         string  `json:"type"`
	Created      float32 `json:"created"`
	Modified     float32 `json:"modified"`
	Model        string  `json:"model"`
	Pushable     bool    `json:"pushable"`
}

//Contact describes a contact entry.
type Contact struct {
	ID              string  `json:"iden"`
	Name            string  `json:"name"`
	Created         float32 `json:"created"`
	Modified        float32 `json:"modified"`
	Email           string  `json:"email"`
	EmailNormalized string  `json:"email_normalized"`
	Active          bool    `json:"active"`
}

//Subscription describes a channel subscription.
type Subscription struct {
	ID       string  `json:"iden"`
	Created  float32 `json:"created"`
	Modified float32 `json:"modified"`
	Active   bool    `json:"active"`
	Channel  Channel `json:"channel"`
}

//Channel describes a channel on a subscription.
type Channel struct {
	ID          string `json:"iden"`
	Tag         string `json:"tag"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}
