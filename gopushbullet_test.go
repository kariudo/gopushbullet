package pushbullet

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func getKey() (k string, err error) {
	k = os.Getenv("APIKEY_PUSHBULLET")
	if len(k) == 0 {
		return k, errors.New("API key env var was not found")
	}
	return
}

func mockHTTP(status int, body string) (*httptest.Server, *Client) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		w.Header().Add("content-type", "application/json")
		fmt.Fprintln(w, body)
	}))

	tr := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}
	httpClient := &http.Client{Transport: tr}

	client := &Client{"apikey", server.URL, httpClient}
	return server, client
}

func TestGetUser(t *testing.T) {
	k, err := getKey()
	if err != nil {
		t.Fatal(err)
	}
	c := ClientWithKey(k)
	u, err := c.GetUser()
	if err != nil {
		t.Error(err)
	}
	p, err := json.MarshalIndent(u, "", "  ")
	fmt.Println(string(p))
}

func TestErrorString(t *testing.T) {
	e := &Error{
		ErrorBody: errorBody{
			Message: "Test invalid request",
			Type:    "invalid_request",
			Cat:     "^._.^",
		},
	}
	if e.String() != "Invalid Request: Test invalid request" {
		t.Error("Did not return correct invalid request error string\nReturned:", e.String())
	}
	e = &Error{
		ErrorBody: errorBody{
			Message: "Test server error",
			Type:    "server",
			Cat:     "^._.^",
		},
	}
	if e.String() != "Server Error: Test server error" {
		t.Error("Did not return correct server error string\nReturned:", e.String())
	}
}

// Devices tests
func TestGetDevices(t *testing.T) {
	k, err := getKey()
	if err != nil {
		t.Error("Failed to get key")
	}
	c := ClientWithKey(k)
	d, err := c.GetDevices()
	if err != nil {
		t.Error("Failed to get devices: ", err)
	}
	fmt.Println(d)
}

// Push - Notes

func TestSendNoteToAll(t *testing.T) {
	// Use the following code in place of the mock calls to test on live api
	// k, err := getKey()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	//c := ClientWithKey(k)
	mockServer, c := mockHTTP(200, "{}")
	defer mockServer.Close()

	err := c.SendNote("Build Test", "This is a test of gopushbullet's SendNote() function.")
	if err != nil {
		t.Error(err)
	}
}

func TestSendNoteFailurePaths(t *testing.T) {
	mockServer, c := mockHTTP(401, "{}")
	defer mockServer.Close()

	err := c.SendNoteToTarget("channel", "testchannelpleaseignore", "Build Test", "This is a test of gopushbullet.")
	if err == nil {
		t.Error(err)
	}
	mockServer, c = mockHTTP(401, "invalid json")
	err = c.SendNoteToTarget("channel", "testchannelpleaseignore", "Build Test", "This is a test of gopushbullet.")
	if err == nil {
		t.Error(err)
	}
}

func TestSendNoteToDevice(t *testing.T) {
	mockServer, c := mockHTTP(200, "{}")
	defer mockServer.Close()

	err := c.SendNoteToTarget("device", "_deviceid_", "Build Test", "This is a test of gopushbullet.")
	if err != nil {
		t.Error(err)
	}
}

func TestSendNoteInvalidTarget(t *testing.T) {
	mockServer, c := mockHTTP(200, "{}")
	defer mockServer.Close()

	err := c.SendNoteToTarget("waffles", "bacon", "Build Test", "This is a test of gopushbullet.")
	if err == nil {
		t.Error(err)
	}
}

func TestSendNoteToChannel(t *testing.T) {
	mockServer, c := mockHTTP(200, "{}")
	defer mockServer.Close()

	err := c.SendNoteToTarget("channel", "testchannelpleaseignore", "Build Test", "This is a test of gopushbullet.")
	if err != nil {
		t.Error(err)
	}
}

func TestSendNoteToEmail(t *testing.T) {
	mockServer, c := mockHTTP(200, "{}")
	defer mockServer.Close()

	err := c.SendNoteToTarget("email", "kariudo@gmail.com", "Build Test", "This is a test of gopushbullet.")
	if err != nil {
		t.Error(err)
	}
}

func TestSendNoteToClientID(t *testing.T) {
	mockServer, c := mockHTTP(200, "{}")
	defer mockServer.Close()

	err := c.SendNoteToTarget("client", "_clientid_", "Build Test", "This is a test of gopushbullet's SendNote() function.")
	if err != nil {
		t.Error(err)
	}
}

// Push - Links
func TestSendLinkToAll(t *testing.T) {
	mockServer, c := mockHTTP(200, "{}")
	defer mockServer.Close()

	err := c.SendLink("Build Test", "This is a test of gopushbullet's SendLink() function.", "http://example.com")
	if err != nil {
		t.Error(err)
	}
}

func TestSendLinkFailurePaths(t *testing.T) {
	mockServer, c := mockHTTP(401, "{}")
	defer mockServer.Close()

	err := c.SendLinkToTarget("channel", "testchannelpleaseignore", "Build Test", "This is a test of gopushbullet's SendLink() function.", "http://example.com")
	if err == nil {
		t.Error(err)
	}
	mockServer, c = mockHTTP(401, "invalid json")
	err = c.SendLinkToTarget("channel", "testchannelpleaseignore", "Build Test", "This is a test of gopushbullet's SendLink() function.", "http://example.com")
	if err == nil {
		t.Error(err)
	}
}

func TestSendLinkToDevice(t *testing.T) {
	mockServer, c := mockHTTP(200, "{}")
	defer mockServer.Close()

	err := c.SendLinkToTarget("device", "_deviceid_", "Build Test", "This is a test of gopushbullet's SendLink() function.", "http://example.com")
	if err != nil {
		t.Error(err)
	}
}

func TestSendLinkInvalidTarget(t *testing.T) {
	mockServer, c := mockHTTP(200, "{}")
	defer mockServer.Close()

	err := c.SendLinkToTarget("waffles", "bacon", "Build Test", "This is a test of gopushbullet's SendLink() function.", "http://example.com")
	if err == nil {
		t.Error(err)
	}
}

func TestSendLinkToChannel(t *testing.T) {
	mockServer, c := mockHTTP(200, "{}")
	defer mockServer.Close()

	err := c.SendLinkToTarget("channel", "testchannelpleaseignore", "Build Test", "This is a test of gopushbullet's SendLink() function.", "http://example.com")
	if err != nil {
		t.Error(err)
	}
}

func TestSendLinkToEmail(t *testing.T) {
	mockServer, c := mockHTTP(200, "{}")
	defer mockServer.Close()

	err := c.SendLinkToTarget("email", "kariudo@gmail.com", "Build Test", "This is a test of gopushbullet's SendLink() function.", "http://example.com")
	if err != nil {
		t.Error(err)
	}
}

func TestSendLinkToClientID(t *testing.T) {
	mockServer, c := mockHTTP(200, "{}")
	defer mockServer.Close()

	err := c.SendLinkToTarget("client", "_clientid_", "Build Test", "This is a test of gopushbullet's SendLink() function.", "http://example.com")
	if err != nil {
		t.Error(err)
	}
}

// Push - Address
func TestSendAddressToAll(t *testing.T) {
	mockServer, c := mockHTTP(200, "{}")
	defer mockServer.Close()

	err := c.SendAddress("Build Test", "Place", "123 Main st., Newtown, CT")
	if err != nil {
		t.Error(err)
	}
}

func TestSendAddressFailurePaths(t *testing.T) {
	mockServer, c := mockHTTP(401, "{}")
	defer mockServer.Close()

	err := c.SendAddressToTarget("channel", "testchannelpleaseignore", "Build Test", "Place", "123 Main st., Newtown, CT")
	if err == nil {
		t.Error(err)
	}
	mockServer, c = mockHTTP(401, "invalid json")
	err = c.SendAddressToTarget("channel", "testchannelpleaseignore", "Build Test", "Place", "123 Main st., Newtown, CT")
	if err == nil {
		t.Error(err)
	}
}

func TestSendAddressToDevice(t *testing.T) {
	mockServer, c := mockHTTP(200, "{}")
	defer mockServer.Close()

	err := c.SendAddressToTarget("device", "_deviceid_", "Build Test", "Place", "123 Main st., Newtown, CT")
	if err != nil {
		t.Error(err)
	}
}

func TestSendAddressInvalidTarget(t *testing.T) {
	mockServer, c := mockHTTP(200, "{}")
	defer mockServer.Close()

	err := c.SendAddressToTarget("waffles", "bacon", "Build Test", "Place", "123 Main st., Newtown, CT")
	if err == nil {
		t.Error(err)
	}
}

func TestSendAddressToChannel(t *testing.T) {
	mockServer, c := mockHTTP(200, "{}")
	defer mockServer.Close()

	err := c.SendAddressToTarget("channel", "testchannelpleaseignore", "Build Test", "Place", "123 Main st., Newtown, CT")
	if err != nil {
		t.Error(err)
	}
}

func TestSendAddressToEmail(t *testing.T) {
	mockServer, c := mockHTTP(200, "{}")
	defer mockServer.Close()

	err := c.SendAddressToTarget("email", "kariudo@gmail.com", "Build Test", "Place", "123 Main st., Newtown, CT")
	if err != nil {
		t.Error(err)
	}
}

func TestSendAddressToClientID(t *testing.T) {
	mockServer, c := mockHTTP(200, "{}")
	defer mockServer.Close()

	err := c.SendAddressToTarget("client", "_clientid_", "Build Test", "Place", "123 Main st., Newtown, CT")

	if err != nil {
		t.Error(err)
	}
}

// Push - Checklist
func TestSendChecklistToAll(t *testing.T) {
	mockServer, c := mockHTTP(200, "{}")
	defer mockServer.Close()

	err := c.SendChecklist("Build Test", []string{"item1", "item2", "item3"})
	if err != nil {
		t.Error(err)
	}
}

func TestSendChecklistFailurePaths(t *testing.T) {
	mockServer, c := mockHTTP(401, "{}")
	defer mockServer.Close()

	err := c.SendChecklistToTarget("channel", "testchannelpleaseignore", "Build Test", []string{"item1", "item2", "item3"})
	if err == nil {
		t.Error(err)
	}
	mockServer, c = mockHTTP(401, "invalid json")
	err = c.SendChecklistToTarget("channel", "testchannelpleaseignore", "Build Test", []string{"item1", "item2", "item3"})
	if err == nil {
		t.Error(err)
	}
}

func TestSendChecklistToDevice(t *testing.T) {
	mockServer, c := mockHTTP(200, "{}")
	defer mockServer.Close()

	err := c.SendChecklistToTarget("device", "_deviceid_", "Build Test", []string{"item1", "item2", "item3"})
	if err != nil {
		t.Error(err)
	}
}

func TestSendChecklistInvalidTarget(t *testing.T) {
	mockServer, c := mockHTTP(200, "{}")
	defer mockServer.Close()

	err := c.SendChecklistToTarget("waffles", "bacon", "Build Test", []string{"item1", "item2", "item3"})
	if err == nil {
		t.Error(err)
	}
}

func TestSendChecklistToChannel(t *testing.T) {
	mockServer, c := mockHTTP(200, "{}")
	defer mockServer.Close()

	err := c.SendChecklistToTarget("channel", "testchannelpleaseignore", "Build Test", []string{"item1", "item2", "item3"})
	if err != nil {
		t.Error(err)
	}
}

func TestSendChecklistToEmail(t *testing.T) {
	mockServer, c := mockHTTP(200, "{}")
	defer mockServer.Close()

	err := c.SendChecklistToTarget("email", "kariudo@gmail.com", "Build Test", []string{"item1", "item2", "item3"})
	if err != nil {
		t.Error(err)
	}
}

func TestSendChecklistToClientID(t *testing.T) {
	mockServer, c := mockHTTP(200, "{}")
	defer mockServer.Close()

	err := c.SendChecklistToTarget("client", "_clientid_", "Build Test", []string{"item1", "item2", "item3"})

	if err != nil {
		t.Error(err)
	}
}

// Contacts
func TestGetContacts(t *testing.T) {
	contactJson := `{
			"contacts": [
				{
				"iden": "ubdcjAfszs0Smi",
				"name": "Ryan Oldenburg",
				"created": 13990116604298899,
				"modified": 139901166042976,
				"email": "ryanjoldenburg@gmail.com",
				"email_normalized": "ryanjoldenburg@gmail.com",
				"active": true
				}
			]
		}`
	mockServer, c := mockHTTP(200, contactJson)
	defer mockServer.Close()

	contacts, err := c.GetContacts()
	if err != nil {
		t.Error(contacts, err)
	}

	if contacts.Contacts[0].Name != "Ryan Oldenburg" {
		t.Error("Contact name not as expected:", contacts.Contacts[0].Name)
	}

}

// Update and Delete Push
func TestDeletePush(t *testing.T) {
	mockServer, c := mockHTTP(200, "{}")
	defer mockServer.Close()
	err := c.DeletePush("pushid")
	if err != nil {
		t.Error("Failure calling DeletePush:", err)
	}
}
