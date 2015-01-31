package pushbullet

import (
	"encoding/json"
	"errors"
	"fmt"
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
