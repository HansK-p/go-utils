package utils

import (
	log "github.com/sirupsen/logrus"
)

// NewLogger will create and return a Logrus base logger
func NewLogger() *log.Logger {
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
