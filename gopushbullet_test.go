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

func TestSendNoteToAll(t *testing.T) {
	// Use the following code in place of the mock calls to test on live api
	// k, err := getKey()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	//c := ClientWithKey(k)
	mockServer, c := mockHTTP(200, "")
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
	mockServer, c := mockHTTP(200, "")
	defer mockServer.Close()

	err := c.SendNoteToTarget("email", "kariudo@gmail.com", "Build Test", "This is a test of gopushbullet.")
	if err != nil {
		t.Error(err)
	}
}

func TestSendNoteToClientID(t *testing.T) {
	mockServer, c := mockHTTP(200, "")
	defer mockServer.Close()

	err := c.SendNoteToTarget("client", "_clientid_", "Build Test", "This is a test of gopushbullet's SendNote() function.")
	if err != nil {
		t.Error(err)
	}
}
