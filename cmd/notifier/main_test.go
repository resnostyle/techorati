package main_test

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"

    "github.com/resnostyle/techorati/pkg/parser"
    "io/ioutil"
    "strings"
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

    t.Run("Successful notification", func(t *testing.T) {
        // ... existing test case ...
    })

    t.Run("Error response", func(t *testing.T) {
        errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.WriteHeader(http.StatusBadRequest)
            w.Write([]byte(`{"status":0,"request":"invalid-request-id"}`))
        }))
        defer errorServer.Close()

        pushover := &main.Pushover{
            User: "testUser",
            Key:  "testKey",
            URL:  errorServer.URL,
        }

        err := pushover.SendNotification("testMessage")
        assert.Error(t, err, "expected an error")
        assert.Contains(t, err.Error(), "invalid-request-id", "error should contain the request ID")
    })

    t.Run("Verify request body", func(t *testing.T) {
        bodyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            body, err := ioutil.ReadAll(r.Body)
            assert.NoError(t, err, "reading body should not error")
            assert.Contains(t, string(body), "token=testKey", "body should contain token")
            assert.Contains(t, string(body), "user=testUser", "body should contain user")
            assert.Contains(t, string(body), "message=testMessage", "body should contain message")
            w.WriteHeader(http.StatusOK)
        }))
        defer bodyServer.Close()

        // ... create pushover and send notification ...
    })

    t.Run("Long message", func(t *testing.T) {
        longMessage := strings.Repeat("a", 1024) // Pushover limit is 1024 characters
        // ... create server, pushover, and send notification with longMessage ...
        // Assert that the message is truncated or an error is returned
    })

    t.Run("Rate limiting", func(t *testing.T) {
        rateLimitServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("X-Limit-App-Remaining", "0")
            w.WriteHeader(http.StatusTooManyRequests)
        }))
        defer rateLimitServer.Close()

        // ... create pushover and send notification ...
        // Assert that an appropriate error is returned
    })
}
