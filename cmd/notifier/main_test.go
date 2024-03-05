package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestHandlePostToSlack(t *testing.T) {
    // Create a new HTTP request
    reqBody := []byte(`{"message":"test message"}`)
    req, err := http.NewRequest("POST", "/post-to-slack", bytes.NewBuffer(reqBody))
    if err != nil {
        t.Fatal(err)
    }

    // Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
    rr := httptest.NewRecorder()

    // Call the handler function directly and pass in the request and response recorder
    handlePostToSlack(rr, req)

    // Check the status code
    assert.Equal(t, http.StatusOK, rr.Code)

    // Check the response body
    expected := []byte("Message posted to Slack successfully")
    assert.Equal(t, expected, rr.Body.Bytes())
}

func TestPostToSlack(t *testing.T) {
    // Mock HTTP server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Check request payload
        var data map[string]interface{}
        err := json.NewDecoder(r.Body).Decode(&data)
        if err != nil {
            t.Fatal(err)
        }

        assert.Equal(t, "test message", data["text"])

        // Respond with a mock success response
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"ok":true}`))
    }))
    defer server.Close()

    // Call the function to be tested
    err := postToSlack("test message", server.URL)
    assert.NoError(t, err)
}
