package main

import (
    "fmt"
    "log"
    "net/http"
    "net/url"
    "os"

    "your-project/envconfig"

    mqtt "github.com/eclipse/paho.mqtt.golang"
)

// Config represents the configuration struct for your service.
type Config struct {
    MQTTBroker   string `env:"MQTT_BROKER,default=tcp://localhost:1883"`
    MQTTTopic    string `env:"MQTT_TOPIC,default=test"`
    PushoverUser string `env:"PUSHOVER_USER"`
    PushoverKey  string `env:"PUSHOVER_KEY"`
    PushoverURL  string `env:"PUSHOVER_URL,default=https://api.pushover.net/1/messages.json"`
}

// Pushover represents a Pushover notification client.
type Pushover struct {
    User string
    Key  string
    URL  string
}

// SendNotification sends a notification to Pushover.
func (p *Pushover) SendNotification(message string) error {
    form := url.Values{}
    form.Add("token", p.Key)
    form.Add("user", p.User)
    form.Add("message", message)

    resp, err := http.PostForm(p.URL, form)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    return nil
}

// MessageHandler is the callback function that handles incoming MQTT messages.
func MessageHandler(client mqtt.Client, msg mqtt.Message) {
    config := loadConfig()

    pushover := &Pushover{
        User: config.PushoverUser,
        Key:  config.PushoverKey,
        URL:  config.PushoverURL,
    }

    err := pushover.SendNotification(string(msg.Payload()))
    if err != nil {
        log.Printf("Error sending notification: %v", err)
    }
}

func loadConfig() *Config {
    var config Config
    err := envconfig.LoadConfig(&config)
    if err != nil {
        log.Fatalf("Error loading configuration: %v", err)
    }
    return &config
}

func main() {
    config := loadConfig()

    opts := mqtt.NewClientOptions().AddBroker(config.MQTTBroker)
    opts.SetDefaultPublishHandler(MessageHandler)

    client := mqtt.NewClient(opts)
    if token := client.Connect(); token.Wait() && token.Error() != nil {
        panic(token.Error())
    }

    topic := fmt.Sprintf("%s/#", config.MQTTTopic)
    if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
        fmt.Println(token.Error())
        os.Exit(1)
    }

    fmt.Printf("Subscribed to MQTT topic: %s\n", topic)
    blocker := make(chan struct{})
    <-blocker
}