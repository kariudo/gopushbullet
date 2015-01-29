package gopushbullet

// ErrorResponse (any non-200 error code) contain information on the kind of error that happened.
type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Cat     string `json:"cat"`
	} `json:"error"`
}
