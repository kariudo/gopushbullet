package gopushbullet

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

//ErrorResponse (any non-200 error code) contain information on the kind of error that happened.
type ErrorResponse struct {
	ErrorBody struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Cat     string `json:"cat"`
	} `json:"error"`
}

func (e ErrorResponse) Error() string {
	var t string
	if e.ErrorBody.Type == "invalid_request" {
		t = "Invalid Request"
	} else {
		t = "Server Error"
	}
	return fmt.Sprintf("%v: %v", e.ErrorBody.Message, t)
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

//User describes the authenticated user.
type User struct {
	ID              string      `json:"iden"`
	Email           string      `json:"email"`
	EmailNormalized string      `json:"email_normalized"`
	Created         float32     `json:"created"`
	Modified        float32     `json:"modified"`
	Name            string      `json:"name"`
	ImageURL        string      `json:"image_url"`
	Preferences     Preferences `json:"preferences"`
}

//Preferences describes a set of user perferences.
type Preferences struct {
	Onboarding struct {
		App       bool `json:"app"`
		Friends   bool `json:"friends"`
		Extension bool `json:"extension"`
	} `json:"onboarding"`
	Social bool   `json:"social"`
	Cat    string `json:"cat"`
}

const baseURL = "https://api.pushbullet.com/v2/"

//GetUser gets the current authenticate users details.
func GetUser(key string) (User, error) {
	var u User
	if len(key) == 0 {
		return u, errors.New("Error: API key required.")
	}
	r, err := makeCall(key, "GET", "users/me", nil)
	if err != nil {
		return u, err
	}
	err = json.Unmarshal(r, &u)
	if err != nil {
		return u, err
	}
	return u, nil
}

func makeCall(key string, method string, call string, body []byte) ([]byte, error) {
	client := &http.Client{}
	r, err := http.NewRequest(method, baseURL+call, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	r.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(key+":")))
	r.Header.Add("Content-Type", "application/json")
	s, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	defer s.Body.Close()

	body, err = ioutil.ReadAll(s.Body)
	if err != nil {
		return nil, err
	}

	if s.StatusCode != http.StatusOK {
		var errResponse ErrorResponse
		err = json.Unmarshal(body, &errResponse)
		if err != nil {
			panic(err)
		}
		return body, &errResponse
	}

	return body, nil
}
