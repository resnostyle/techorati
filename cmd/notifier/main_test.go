package main_test

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"

    "github.com/resnostyle/techorati/pkg/parser"
)

func TestSendNotification(t *testing.T) {
    // Create a test server to mock the Pushover API
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        assert.Equal(t, http.MethodPost, r.Method, "expected POST method")
        assert.Equal(t, "/1/messages.json", r.URL.String(), "expected URL /1/messages.json")
        assert.Equal(t, "testKey", r.PostFormValue("token"), "expected token testKey")
        assert.Equal(t, "testUser", r.PostFormValue("user"), "expected user testUser")
        assert.Equal(t, "testMessage", r.PostFormValue("message"), "expected message testMessage")

        // Respond with success
        w.WriteHeader(http.StatusOK)
    }))
    defer server.Close()

    // Replace the Pushover URL with the test server URL
    pushover := &main.Pushover{
        User: "testUser",
        Key:  "testKey",
        URL:  server.URL,
    }

    // Call the function to test
    err := pushover.SendNotification("testMessage")
    assert.NoError(t, err, "expected no error")
}
