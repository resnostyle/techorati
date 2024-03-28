package envconfig

// Config represents a configuration struct that can be loaded from environment variables.
type Config struct {
	MQTTBroker   string `env:"MQTT_BROKER,default=tcp://localhost:1883"`
	MQTTTopic    string `env:"MQTT_TOPIC,default=test"`
	PushoverUser string `env:"PUSHOVER_USER"`
	PushoverKey  string `env:"PUSHOVER_KEY"`
	PushoverURL  string `env:"PUSHOVER_URL,default=https://api.pushover.net/1/messages.json"`
	SlackURL     string `env:"SLACK_URL"`
}