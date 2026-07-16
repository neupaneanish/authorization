package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Env struct {
	ValkeyURL    string
	ServiceName  string
	TelemetryURL string
	Environment  string
	Port         string
}

const (
	Development = "development"
	Production  = "production"
)

func LoadEnv() (*Env, error) {
	valkeyURL, valkeyURLErr := ValidateEnv("VALKEY_URL")
	if valkeyURLErr != nil {
		return nil, valkeyURLErr
	}

	environment := ValidateDefaultEnv("ENVIRONMENT", Development)
	switch environment {
	case Development, Production:
	default:
		return nil, fmt.Errorf("ENVIRONMENT must be %s or %s", Development, Production)
	}

	telemetryURL, telemetryURLErr := ValidateEnv("TELEMETRY_URL")
	if telemetryURLErr != nil {
		return nil, telemetryURLErr
	}

	port := ValidateDefaultEnv("PORT", "50051")
	value, valueErr := strconv.Atoi(port)
	if valueErr != nil || value < 80 || value > 65535 {
		return nil, errors.New("PORT must be between 80  and 65535")
	}

	return &Env{
		ValkeyURL:    valkeyURL,
		TelemetryURL: telemetryURL,
		ServiceName:  ValidateDefaultEnv("SERVICE_NAME", "neupaneanish.com.np/authorization"),
		Environment:  environment,
		Port:         port,
	}, nil
}

func ValidateEnv(key string) (string, error) {
	env := os.Getenv(key)
	value := strings.TrimSpace(env)
	if value == "" {
		return "", fmt.Errorf("%s is missing", key)
	}
	return value, nil
}

func ValidateDefaultEnv(key string, def string) string {
	env := os.Getenv(key)
	value := strings.TrimSpace(env)
	if value == "" {
		return def
	}
	return value
}
