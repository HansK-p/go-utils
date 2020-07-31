package utils

import (
	"fmt"
	"net"
	"net/http"

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

// HealthObject interface implement a metod to validate if the object is healthy
type HealthObject interface {
	IsHealthy() (bool, string)
}

// HealthHandler contains what is needed in order to validate health
type HealthHandler struct {
	HealthObject HealthObject
}

func (hh *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	healthy, message := hh.HealthObject.IsHealthy()
	if healthy {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(message))
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(message))
}

const healthPort = 8080

//RunHTTPHealthListener will start a listener listening for health and liveness checks
func RunHTTPHealthListener(logger *log.Logger, hh *HealthHandler) {
	m := http.NewServeMux()
	m.Handle("/healthz", hh)
	logger.Infof("Starting /healthz endpoint on 0.0.0.0:%d", healthPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", healthPort))
	if err != nil {
		logger.Fatalf("Failed to start Health endpoint: %v", err)
	}
	go http.Serve(lis, m)
}
