package main

import (
    "encoding/json"
    "log"
    "net/http"
    "bytes"
)

// Define a struct to represent the JSON data you expect to receive
type Data struct {
    Message string `json:"message"`
}

func main() {
    http.HandleFunc("/post-to-slack", handlePostToSlack)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func handlePostToSlack(w http.ResponseWriter, r *http.Request) {
    // Decode the JSON data from the request body
    var data Data
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Post the message to Slack
    if err := postToSlack(data.Message); err != nil {
        http.Error(w, "Failed to post message to Slack", http.StatusInternalServerError)
        return
    }

    // Respond with a success message
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Message posted to Slack successfully"))
}

func postToSlack(message string) error {
    // TODO: Implement your code to post message to Slack
    // You'll need to use the Slack API here
    // You can use libraries like slack-go/slack to interact with the Slack API
    // Example: https://github.com/slack-go/slack
    // Here's a simple example using HTTP client:

    // Replace "YOUR_SLACK_API_TOKEN" and "YOUR_CHANNEL_ID" with your actual Slack API token and channel ID
    slackURL := "https://slack.com/api/chat.postMessage"
    slackToken := "YOUR_SLACK_API_TOKEN"
    channelID := "YOUR_CHANNEL_ID"

    // Construct the message payload
    payload := map[string]interface{}{
        "channel": channelID,
        "text":    message,
    }

    // Marshal payload to JSON
    jsonPayload, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    // Send POST request to Slack API
    client := &http.Client{}
    req, err := http.NewRequest("POST", slackURL, bytes.NewBuffer(jsonPayload))
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+slackToken)

    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("Slack API returned non-200 status code: %d", resp.StatusCode)
    }

    return nil
}
