package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    mqtt "github.com/eclipse/paho.mqtt.golang"
    "github.com/gorilla/mux"
    "github.com/kelseyhightower/envconfig"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// Config represents the configuration struct for the API service.
type Config struct {
    MQTTBroker string `env:"MQTT_BROKER,default=tcp://localhost:1883"`
    MQTTTopic  string `env:"MQTT_TOPIC,default=test"`
    APIPort    string `env:"API_PORT,default=8080"`
}

// Message represents the structure of the incoming JSON payload.
type Message struct {
    Content string `json:"content"`
}

var (
    mqttClient mqtt.Client
    logger     *log.Logger

    // Prometheus metrics
    httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
        Name: "http_requests_total",
        Help: "Total number of HTTP requests",
    }, []string{"method", "endpoint", "status"})

    mqttMessagesPublished = promauto.NewCounter(prometheus.CounterOpts{
        Name: "mqtt_messages_published_total",
        Help: "Total number of MQTT messages published",
    })

    mqttConnectionStatus = promauto.NewGauge(prometheus.GaugeOpts{
        Name: "mqtt_connection_status",
        Help: "Status of MQTT connection (1 for connected, 0 for disconnected)",
    })
)

func init() {
    logger = log.New(os.Stdout, "API: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
    config := loadConfig()

    // Set up MQTT client
    opts := mqtt.NewClientOptions().AddBroker(config.MQTTBroker)
    opts.SetClientID("api-service")
    opts.SetConnectTimeout(5 * time.Second)
    opts.SetAutoReconnect(true)
    opts.SetOnConnectHandler(func(client mqtt.Client) {
        logger.Println("Connected to MQTT broker")
        mqttConnectionStatus.Set(1)
    })
    opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
        logger.Printf("Connection to MQTT broker lost: %v", err)
        mqttConnectionStatus.Set(0)
    })

    mqttClient = mqtt.NewClient(opts)
    if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
        logger.Fatalf("Error connecting to MQTT broker: %v", token.Error())
    }

    // Set up HTTP router
    r := mux.NewRouter()
    r.HandleFunc("/send", handleSendMessage).Methods("POST")
    r.HandleFunc("/health", handleHealthCheck).Methods("GET")
    r.Handle("/metrics", promhttp.Handler())

    // Wrap the router with the Prometheus middleware
    promRouter := prometheusMiddleware(r)

    // Start the server
    logger.Printf("Starting server on port %s\n", config.APIPort)
    server := &http.Server{
        Addr:         ":" + config.APIPort,
        Handler:      promRouter,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
    }
    log.Fatal(server.ListenAndServe())
}

func prometheusMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        route := mux.CurrentRoute(r)
        path, _ := route.GetPathTemplate()
        
        rw := NewResponseWriter(w)
        next.ServeHTTP(rw, r)
        
        statusCode := rw.statusCode
        httpRequestsTotal.WithLabelValues(r.Method, path, fmt.Sprintf("%d", statusCode)).Inc()
    })
}

type ResponseWriter struct {
    http.ResponseWriter
    statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
    return &ResponseWriter{w, http.StatusOK}
}

func (rw *ResponseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func handleSendMessage(w http.ResponseWriter, r *http.Request) {
    var msg Message
    if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
        logger.Printf("Error decoding request body: %v", err)
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    if msg.Content == "" {
        http.Error(w, "Message content cannot be empty", http.StatusBadRequest)
        return
    }

    config := loadConfig()
    token := mqttClient.Publish(config.MQTTTopic, 0, false, msg.Content)
    token.Wait()

    if token.Error() != nil {
        logger.Printf("Error publishing message: %v", token.Error())
        http.Error(w, "Failed to publish message", http.StatusInternalServerError)
        return
    }

    mqttMessagesPublished.Inc()
    logger.Printf("Message published: %s", msg.Content)
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "Message sent"})
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
    if !mqttClient.IsConnected() {
        w.WriteHeader(http.StatusServiceUnavailable)
        json.NewEncoder(w).Encode(map[string]string{"status": "unhealthy", "reason": "MQTT client disconnected"})
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func loadConfig() *Config {
    var config Config
    err := envconfig.Process("", &config)
    if err != nil {
        logger.Fatalf("Error loading configuration: %v", err)
    }
    return &config
}