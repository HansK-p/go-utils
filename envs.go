package utils

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// Getenv will return OS environment with default value. Will exit if a mandatory envrionment variable isn't set
func Getenv(logger *log.Logger, envName, defaultValue string, mandatory bool) string {
	value := os.Getenv(envName)
	if value != "" {
		return value
	}
	if defaultValue == "" && mandatory {
		logger.Fatalf("The environment variable '%s' must be set", envName)
	}
	return defaultValue
}
