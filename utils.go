package utils

import (
	log "github.com/sirupsen/logrus"

	"os"
)

// GetLogger will create and return a Logrus base logger
func GetLogger() *log.Logger {
	logger := log.New()
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	loglevel, err := log.ParseLevel(Getenv(logger, "LOG_LEVEL", "INFO", false))
	if err != nil {
		logger.Fatalf("Unable to parse environment variable LOG_LEVEL: %v", err)
	}
	logger.SetLevel(loglevel)
	return logger
}

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
