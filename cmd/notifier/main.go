package main

import (
    "fmt"
    "log"
    "net/http"
    "net/url"
    "os"

    "github.com/resnostyle/techorati/pkg/parser"

    mqtt "github.com/eclipse/paho.mqtt.golang"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
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
    timer := prometheus.NewTimer(notificationDuration)
    defer timer.ObserveDuration()

    form := url.Values{}
    form.Add("token", p.Key)
    form.Add("user", p.User)
    form.Add("message", message)

    resp, err := http.PostForm(p.URL, form)
    if err != nil {
        notificationErrors.Inc()
        return err
    }
    defer resp.Body.Close()

    notificationsSent.Inc()
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

    if err := pushover.SendNotification(string(msg.Payload())); err != nil {
        log.Printf("Error sending notification: %v", err)
    }
}

func init() {
    logger = log.New(os.Stdout, "API: ", log.Ldate|log.Ltime|log.Lshortfile)
    if err := envconfig.Process("", &config); err != nil {
        logger.Fatalf("Error loading configuration: %v", err)
    }
}

func setupMQTTClient() error {
    opts := mqtt.NewClientOptions().AddBroker(config.MQTTBroker)
    // ... (rest of the MQTT setup)
    mqttClient = mqtt.NewClient(opts)
    if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
        return fmt.Errorf("Error connecting to MQTT broker: %w", token.Error())
    }
    return nil
}
}

func main() {
    config := loadConfig()

    // Set up Prometheus HTTP server
    go func() {
        http.Handle("/metrics", promhttp.Handler())
        log.Fatal(http.ListenAndServe(":2112", nil))
    }()

    if err := setupMQTTClient(); err != nil {
        logger.Fatalf("Failed to set up MQTT client: %v", err)
    }

    topic := fmt.Sprintf("%s/#", config.MQTTTopic)
    if token := mqttClient.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
        fmt.Println(token.Error())
        os.Exit(1)
    }

    fmt.Printf("Subscribed to MQTT topic: %s\n", topic)
    fmt.Println("Prometheus metrics available at :2112/metrics")
    blocker := make(chan struct{})
    <-blocker
}