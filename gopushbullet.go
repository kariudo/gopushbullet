package gopushbullet

// ErrorResponse (any non-200 error code) contain information on the kind of error that happened.
type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Cat     string `json:"cat"`
	} `json:"error"`
}

// PushMessage describes a message to be sent via PushBullet. Only one of the first 4 properties may be specified with a message being sent.
type PushMessage struct {
	DeviceIdentity string `json:"device_iden"`
	Email          string `json:"email"`
	ChannelTag     string `json:"channel_tag"`
	ClientIdentity string `json:"client_iden"`

	// Type indicates the type of push message being sent
	Type    string   `json:"type"`    // note, link, address, checklist, or file
	Title   string   `json:"title"`   // Title is used withn note, link, and checklist types
	Body    string   `json:"body"`    // Body is used with link, note, or file types
	URL     string   `json:"url"`     // URL is used for link types
	Address string   `json:"address"` // Address is used with address type
	Items   []string `json:"items"`   // Items are used with checklist types
	// The following are used with file types
	FileName             string `json:"file_name"`
	FileType             string `json:"file_type"` // MIME type of the file
	FileURL              string `json:"file_url"`
	SourceDeviceIdentity string `json:"source_device_iden"`
}
