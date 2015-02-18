package pushbullet

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
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

//PushMessage describes a message to be sent via Pushbullet. Only one of the first 4 properties may be specified with a message being sent.
type PushMessage struct {
	ID string `json:"`

	// Target specific properties
	DeviceID   string `json:"device_iden"`
	Email      string `json:"email"`
	ChannelTag string `json:"channel_tag"`
	ClientID   string `json:"client_iden"`

	// Type indicates the type of push message being sent
	Type    string   `json:"type"`    // note, link, address, checklist, or file
	Title   string   `json:"title"`   // Title is used within note, link, and checklist types
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

	// Properties for response messages
	Created                 float32 `json:"created"`
	Modified                float32 `json:"modified"`
	Active                  bool    `json:"active"`
	Dismissed               bool    `json:"dismissed"`
	SenderID                string  `json:"sender_iden"`
	SenderEmail             string  `json:"sender_email"`
	SenderEmailNormalized   string  `json:"sender_email_normalized"`
	ReceiverID              string  `json:"receiver_iden"`
	ReceiverEmail           string  `json:"receiver_email"`
	ReceiverEmailNormalized string  `json:"receiver_email_normalized"`
}

//PushList describes a list of push messages
type PushList struct {
	Pushes []PushMessage `json:"pushes"`
}

//ItemsList describes a list of checklist items
type ItemsList struct {
	items []Item `json:"items"`
}

//Item describes a checklist item
type Item struct {
	Text    string `json:"text"`
	Checked bool   `json:"checked"`
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

//SubscriptionList describes a list of subscribed channels
type SubscriptionList struct {
	Subscriptions []Subscription `json:"subscriptions"`
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

//Preferences describes a set of user preferences.
type Preferences struct {
	Onboarding struct {
		App       bool `json:"app"`
		Friends   bool `json:"friends"`
		Extension bool `json:"extension"`
	} `json:"onboarding"`
	Social bool   `json:"social"`
	Cat    string `json:"cat"`
}

//Authorization describes a file upload authorization.
type Authorization struct {
	FileType  string `json:"file_type"`
	FileName  string `json:"file_name"`
	FileURL   string `json:"file_url"`
	UploadURL string `json:"upload_url"`
	Data      struct {
		Awsaccesskeyid string `json:"awsaccesskeyid"`
		Acl            string `json:"acl"`
		Key            string `json:"key"`
		Signature      string `json:"signature"`
		Policy         string `json:"policy"`
		ContentType    string `json:"content-type"`
	} `json:"data"`
}

//Client a Pushbullet API client
type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

//ClientWithKey returns a pushbullet.Client pointer with API key.
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
		// only remaining acceptable type is "all" which takes no additional fields
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
		// only remaining acceptable type is "all" which takes no additional fields
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
		// only remaining acceptable type is "all" which takes no additional fields
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
		// only remaining acceptable type is "all" which takes no additional fields
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
		// only remaining acceptable type is "all" which takes no additional fields
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

//UpdateContact creates a new contact with the specified name and email
func (c *Client) UpdateContact(contactID, name string) error {
	u := url.Values{}
	u.Add("name", name)
	_, err := c.HTTPClient.PostForm(c.BaseURL+"contacts/"+contactID, u)
	if err != nil {
		return err
	}
	return nil
}

//DeleteContact creates a new contact with the specified name and email
func (c *Client) DeleteContact(contactID string) error {
	_, apiError, err := c.makeCall("DELETE", "contacts/"+contactID, nil)
	if err != nil {
		log.Println("Failed to delete contact: ", err, apiError.String())
		return err
	}
	return nil
}

//SubscribeChannel subscribes use to a specified channel
func (c *Client) SubscribeChannel(channel string) error {
	_, apiError, err := c.makeCall("POST", "subscriptions", nil)
	if err != nil {
		log.Println("Failed to add subscription: ", err, apiError.String())
		return err
	}
	return nil
}

//ListSubscriptions returns a list of channels to which the user is subscribed
func (c *Client) ListSubscriptions() (subscriptions SubscriptionList, err error) {
	responseBody, apiError, err := c.makeCall("GET", "subscriptions", nil)
	if err != nil {
		log.Println("Failed to add subscription: ", err, apiError.String())
		return
	}
	err = json.Unmarshal(responseBody, &subscriptions)
	if err != nil {
		return
	}
	return
}

//UnsubscribeChannel unsubscribes from the specified channel
func (c *Client) UnsubscribeChannel(channelID string) error {
	_, apiError, err := c.makeCall("DELETE", "subscriptions/"+channelID, nil)
	if err != nil {
		log.Println("Failed to unsubscribe channel: ", err, apiError.String())
		return err
	}
	return nil
}

//ChannelInfo gets detained info for the requested channel
func (c *Client) ChannelInfo(channelTag string) (channel Channel, err error) {
	response, apiError, err := c.makeCall("GET", "channel-info?tag="+channelTag, nil)
	if err != nil {
		log.Println("Failed to get channel info: ", err, apiError.String())
		return
	}
	err = json.Unmarshal(response, &channel)
	return
}

//AuthorizeUpload requests an authorization to upload a file
func (c *Client) AuthorizeUpload(fileName, fileType string) (Authorization, error) {
	var auth Authorization
	u := url.Values{}
	u.Add("file_name", fileName)
	u.Add("file_type", fileType)
	response, err := c.HTTPClient.PostForm(c.BaseURL+"upload-request", u)
	if err != nil {
		return auth, err
	}
	// read the response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return auth, err
	}
	err = json.Unmarshal(body, &auth)
	if err != nil {
		return auth, err
	}
	return auth, nil
}

//UpdatePreferences overwrites user preferences with specified ones
func (c *Client) UpdatePreferences(preferences Preferences) error {
	_, apiError, err := c.makeCall("POST", "users/me", preferences)
	if err != nil {
		log.Println("Failed to update preferences: ", apiError, err)
		return err
	}
	return err
}

//GetPushHistory gets pushes modified after the provided timestamp
func (c *Client) GetPushHistory(modifiedAfter float32) ([]PushMessage, error) {
	var pushList PushList
	responseBody, apiError, err := c.makeCall("GET", "pushes?modified_after="+strconv.FormatFloat(float64(modifiedAfter), 'f', 4, 32), nil)
	if err != nil {
		log.Println("Error getting push history: ", apiError, err)
		return pushList.Pushes, err
	}
	err = json.Unmarshal(responseBody, &pushList)
	if err != nil {
		return pushList.Pushes, err
	}
	return pushList.Pushes, nil
}

//DeletePush deletes a push message
func (c *Client) DeletePush(pushID string) error {
	_, apiError, err := c.makeCall("DELETE", "pushes/"+pushID, nil)
	if err != nil {
		log.Println("Failed to delete push: ", apiError, err)
		return err
	}
	return nil
}

//DismissPush allows for dismissal of a push message
func (c *Client) DismissPush(ID string) error {
	_, apiError, err := c.makeCall("GET", "pushes/"+ID, nil)
	if err != nil {
		log.Println("Failed to dismiss push: ", apiError, err)
		return err
	}
	return nil
}

//UpdateList allows for updating a list type push
func (c *Client) UpdateList(pushID string, list ItemsList) error {
	_, apiError, err := c.makeCall("POST", "pushes/"+pushID, list)
	if err != nil {
		log.Println("Failed to update list: ", apiError, err)
		return err
	}
	return nil
}

//makeCall handles most http transactions under standard methods
func (c *Client) makeCall(method string, call string, data interface{}) (responseBody []byte, apiError *Error, err error) {
	// make sure API key seems OK
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

func uploadFileByPath(authorization Authorization, file string) (err error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	// Add file
	f, err := os.Open(file)
	if err != nil {
		return
	}
	fw, err := w.CreateFormFile("file", file)
	if err != nil {
		return
	}
	if _, err = io.Copy(fw, f); err != nil {
		return
	}
	// Add the other fields
	if fw, err = w.CreateFormField("awsaccesskeyid"); err != nil {
		return
	}
	if _, err = fw.Write([]byte(authorization.Data.Awsaccesskeyid)); err != nil {
		return
	}
	if fw, err = w.CreateFormField("acl"); err != nil {
		return
	}
	if _, err = fw.Write([]byte(authorization.Data.Acl)); err != nil {
		return
	}
	if fw, err = w.CreateFormField("key"); err != nil {
		return
	}
	if _, err = fw.Write([]byte(authorization.Data.Key)); err != nil {
		return
	}
	if fw, err = w.CreateFormField("signature"); err != nil {
		return
	}
	if _, err = fw.Write([]byte(authorization.Data.Signature)); err != nil {
		return
	}
	if fw, err = w.CreateFormField("policy"); err != nil {
		return
	}
	if _, err = fw.Write([]byte(authorization.Data.Policy)); err != nil {
		return
	}
	if fw, err = w.CreateFormField("content-type"); err != nil {
		return
	}
	if _, err = fw.Write([]byte(authorization.Data.ContentType)); err != nil {
		return
	}
	w.Close()

	req, err := http.NewRequest("POST", authorization.UploadURL, &b)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}

	// Check the response
	if res.StatusCode >= 300 {
		err = fmt.Errorf("Bad Status Result: %s", res.Status)
	}

	return err
}
