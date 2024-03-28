package envconfig

import "os"

// LoadConfig loads the configuration from environment variables into the provided config struct.
func LoadConfig(config interface{}) error {
	return loadEnvVars(config)
}

// loadEnvVars loads the environment variables into the provided config struct.
func loadEnvVars(config interface{}) error {
	return ParseEnvVars(config, os.Getenv)
}