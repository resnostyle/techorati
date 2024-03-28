package main

import (
    "testing"

    mqtt "github.com/eclipse/paho.mqtt.golang"
)

type mockPushover struct {
    sentMessages []string
}

func (m *mockPushover) SendNotification(message string) error {
    m.sentMessages = append(m.sentMessages, message)
    return nil
}

func TestMessageHandler(t *testing.T) {
    tests := []struct {
        name     string
        payload  string
        expected []string
    }{
        {
            name:     "Send a single message",
            payload:  "Hello, world!",
            expected: []string{"Hello, world!"},
        },
        {
            name:     "Send multiple messages",
            payload:  "Message 1\nMessage 2\nMessage 3",
            expected: []string{"Message 1", "Message 2", "Message 3"},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            mockClient := mqtt.NewClient(mqtt.NewClientOptions())
            mockPush := &mockPushover{}

            config := &Config{
                PushoverUser: "test-user",
                PushoverKey:  "test-key",
                PushoverURL:  "http://example.com/pushover",
            }

            MessageHandler(mockClient, &mqtt.Message{
                Payload: []byte(tc.payload),
            })

            if len(mockPush.sentMessages) != len(tc.expected) {
                t.Errorf("Expected %d messages, got %d", len(tc.expected), len(mockPush.sentMessages))
            }

            for i, msg := range tc.expected {
                if mockPush.sentMessages[i] != msg {
                    t.Errorf("Expected message %q, got %q", msg, mockPush.sentMessages[i])
                }
            }
        })
    }
}