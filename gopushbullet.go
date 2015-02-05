package pushbullet

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

//Error (any non-200 error code) contain information on the kind of error that happened.
type (
	Error struct {
		ErrorBody errorBody `json:"error"`
	}

	errorBody struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Cat     string `json:"cat"`
	}
)

func (e *Error) String() string {
	var t string
	if e.ErrorBody.Type == "invalid_request" {
		t = "Invalid Request"
	} else {
		t = "Server Error"
	}
	return fmt.Sprintf("%v: %v", t, e.ErrorBody.Message)
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
	Name    string   `json:"name"`    // Name of place used with address type
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

//DeviceList describes an array of devices
type DeviceList struct {
	Devices []Device `json:"devices"`
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

//ContactList describes an array of contacts
type ContactList struct {
	Contacts []Contact `json:"contacts"`
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

//Client a Pushbullet API client
type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

//ClientWithKey returns a pushbullet.CLient pointer with API key.
func ClientWithKey(key string) *Client {
	return &Client{
		APIKey:     key,
		BaseURL:    "https://api.pushbullet.com/v2/",
		HTTPClient: &http.Client{},
	}
}

//GetUser gets the current authenticate users details.
func (c *Client) GetUser() (u User, err error) {
	r, apiError, err := c.makeCall("GET", "users/me", nil)
	if err != nil {
		log.Println("Failed to get user:", err, apiError.String())
		return u, err
	}
	err = json.Unmarshal(r, &u)
	if err != nil {
		return u, err
	}
	return u, nil
}

//SendNote simply sends a note type push to all of the users devices
func (c *Client) SendNote(title, body string) error {
	err := c.SendNoteToTarget("all", "", title, body)
	return err
}

//SendNoteToTarget sends a note type push to a specific device.
func (c *Client) SendNoteToTarget(targetType, target, title, body string) error {
	var p = PushMessage{
		Type:  "note",
		Title: title,
		Body:  body,
	}
	switch targetType {
	case "device":
		p.DeviceID = target
	case "email":
		p.Email = target
	case "channel":
		p.ChannelTag = target
	case "client":
		p.ClientID = target
	default:
		// only remaining acceptable type is "all" which takes no addtional fields
		if targetType != "all" {
			return errors.New("Invalid target type")
		}
	}

	_, apiError, err := c.makeCall("POST", "pushes", p)
	if err != nil {
		log.Println("Failed to send note:", err, apiError.String())
		return err
	}
	return nil
}

//SendLink simply sends a link type push to all of the users devices
func (c *Client) SendLink(title, body, url string) error {
	err := c.SendLinkToTarget("all", "", title, body, url)
	return err
}

//SendLinkToTarget sends a link type push to a specific device.
func (c *Client) SendLinkToTarget(targetType, target, title, body, url string) error {
	var p = PushMessage{
		Type:  "link",
		Title: title,
		Body:  body,
		URL:   url,
	}
	switch targetType {
	case "device":
		p.DeviceID = target
	case "email":
		p.Email = target
	case "channel":
		p.ChannelTag = target
	case "client":
		p.ClientID = target
	default:
		// only remaining acceptable type is "all" which takes no addtional fields
		if targetType != "all" {
			return errors.New("Invalid target type")
		}
	}

	_, apiError, err := c.makeCall("POST", "pushes", p)
	if err != nil {
		log.Println("Failed to get user:", err, apiError.String())
		return err
	}
	return nil
}

//SendAddress simply sends an address type push to all of the users devices
func (c *Client) SendAddress(title, name, address string) error {
	err := c.SendLinkToTarget("all", "", title, name, address)
	return err
}

//SendAddressToTarget sends an address type push to a specific device.
func (c *Client) SendAddressToTarget(targetType, target, title, name, address string) error {
	var p = PushMessage{
		Type:    "address",
		Title:   title,
		Name:    name,
		Address: address,
	}
	switch targetType {
	case "device":
		p.DeviceID = target
	case "email":
		p.Email = target
	case "channel":
		p.ChannelTag = target
	case "client":
		p.ClientID = target
	default:
		// only remaining acceptable type is "all" which takes no addtional fields
		if targetType != "all" {
			return errors.New("Invalid target type")
		}
	}

	_, apiError, err := c.makeCall("POST", "pushes", p)
	if err != nil {
		log.Println("Failed to send address:", err, apiError.String())
		return err
	}
	return nil
}

//SendChecklist simply sends a checklist type push to all of the users devices
func (c *Client) SendChecklist(title string, items []string) error {
	err := c.SendChecklistToTarget("all", "", title, items)
	return err
}

//SendChecklistToTarget sends a checklist type push to a specific device.
func (c *Client) SendChecklistToTarget(targetType, target, title string, items []string) error {
	var p = PushMessage{
		Type:  "checklist",
		Title: title,
		Items: items,
	}
	switch targetType {
	case "device":
		p.DeviceID = target
	case "email":
		p.Email = target
	case "channel":
		p.ChannelTag = target
	case "client":
		p.ClientID = target
	default:
		// only remaining acceptable type is "all" which takes no addtional fields
		if targetType != "all" {
			return errors.New("Invalid target type")
		}
	}

	_, apiError, err := c.makeCall("POST", "pushes", p)
	if err != nil {
		log.Println("Failed to send checklist:", err, apiError.String())
		return err
	}
	return nil
}

//SendFile simply sends a file type push to all of the users devices
func (c *Client) SendFile(title string, items []string) error {
	err := c.SendChecklistToTarget("all", "", title, items)
	return err
}

//SendFileToTarget sends a file type push to a specific device.
func (c *Client) SendFileToTarget(targetType, target, fileName, fileType, fileURL, body string, items []string) error {
	var p = PushMessage{
		Type:     "file",
		FileName: fileName,
		FileType: fileType,
		FileURL:  fileURL,
		Body:     body,
	}
	switch targetType {
	case "device":
		p.DeviceID = target
	case "email":
		p.Email = target
	case "channel":
		p.ChannelTag = target
	case "client":
		p.ClientID = target
	default:
		// only remaining acceptable type is "all" which takes no addtional fields
		if targetType != "all" {
			return errors.New("Invalid target type")
		}
	}

	_, apiError, err := c.makeCall("POST", "pushes", p)
	if err != nil {
		log.Println("Failed to send file: ", err, apiError.String())
		return err
	}
	return nil
}

//GetDevices obtains a list of registered devices from Pushbullet
func (c *Client) GetDevices() (DeviceList, error) {
	var d DeviceList
	res, apiError, err := c.makeCall("GET", "devices", nil)
	if err != nil {
		log.Println("Failed to get devices: ", err, apiError.String())
		return d, err
	}
	err = json.Unmarshal(res, &d)
	if err != nil {
		return d, err
	}
	return d, nil
}

//GetContacts obtains a list of your contacts
func (c *Client) GetContacts() (ContactList, error) {
	var l ContactList
	res, apiError, err := c.makeCall("GET", "contacts", nil)
	if err != nil {
		log.Println("Failed to get contacts: ", err, apiError.String())
		return l, err
	}
	err = json.Unmarshal(res, &l)
	if err != nil {
		return l, err
	}
	return l, err
}

//CreateContact creates a new contact with the specified name and email
func (c *Client) CreateContact(name, email string) error {
	u := url.Values{}
	u.Add("name", name)
	u.Add("email", email)
	_, err := c.HTTPClient.PostForm(c.BaseURL+"contacts", u)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) makeCall(method string, call string, data interface{}) (responseBody []byte, apiError *Error, err error) {
	// make sure API key seems ok
	if len(c.APIKey) == 0 {
		return responseBody, apiError, errors.New("Error: API key required.")
	}

	var payload []byte
	// create the payload
	if data != nil {
		payload, err = json.Marshal(data)
		if err != nil {
			return responseBody, apiError, err
		}
	}

	// make the call
	req, err := http.NewRequest(method, c.BaseURL+call, bytes.NewBuffer(payload))
	if err != nil {
		return responseBody, apiError, err
	}
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.APIKey+":")))
	req.Header.Add("Content-Type", "application/json")
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return responseBody, apiError, err
	}
	defer res.Body.Close()

	// read the response
	responseBody, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return responseBody, apiError, err
	}

	// if the response was an error message
	if res.StatusCode != http.StatusOK {
		apiError = &Error{}
		err = json.Unmarshal(responseBody, &apiError)
		if err != nil {
			return responseBody, apiError, err
		}
		return responseBody, apiError, fmt.Errorf("Status code: %v", res.StatusCode)
	}

	return responseBody, apiError, err
}
