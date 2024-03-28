package envconfig

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// EnvVarFunc is a function that retrieves the value of an environment variable.
type EnvVarFunc func(string) string

// ParseEnvVars populates the fields of the provided config struct with values
// from environment variables. It uses the provided EnvVarFunc to retrieve the
// value of the environment variable.
func ParseEnvVars(config interface{}, getEnvVar EnvVarFunc) error {
	configValue := reflect.ValueOf(config).Elem()
	configType := configValue.Type()

	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		envTag := field.Tag.Get("env")
		if envTag == "" {
			continue
		}

		parts := strings.Split(envTag, ",")
		envVarName := parts[0]
		defaultValue := ""
		if len(parts) > 1 {
			defaultValue = parts[1][len("default="):]
		}

		envVarValue := getEnvVar(envVarName)
		if envVarValue == "" {
			envVarValue = defaultValue
		}

		fieldValue := configValue.Field(i)
		if !fieldValue.CanSet() {
			return fmt.Errorf("cannot set field %s", field.Name)
		}

		switch fieldValue.Kind() {
		case reflect.String:
			fieldValue.SetString(envVarValue)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intValue, err := strconv.ParseInt(envVarValue, 10, 64)
			if err != nil {
				return fmt.Errorf("cannot parse %q as int for field %s: %w", envVarValue, field.Name, err)
			}
			fieldValue.SetInt(intValue)
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(envVarValue)
			if err != nil {
				return fmt.Errorf("cannot parse %q as bool for field %s: %w", envVarValue, field.Name, err)
			}
			fieldValue.SetBool(boolValue)
		default:
			return fmt.Errorf("unsupported field type %s for field %s", fieldValue.Kind(), field.Name)
		}
	}

	return nil
}