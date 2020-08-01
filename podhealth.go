package utils

import (
	"fmt"
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// HealthObject interface implement a metod to validate if the object is healthy
type HealthObject interface {
	IsAlive() (bool, string)
	IsReady() (bool, string)
}

// HealthHandler contains what is needed in order to validate health
type HealthHandler struct {
	HealthObject HealthObject
}

func (hh *HealthHandler) IsAlive(w http.ResponseWriter, r *http.Request) {
	health, message := hh.HealthObject.IsAlive()
	if health {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(message))
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(message))
}
func (hh *HealthHandler) IsReady(w http.ResponseWriter, r *http.Request) {
	health, message := hh.HealthObject.IsReady()
	if health {
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
	m.HandleFunc("/healthz", hh.IsAlive)
	m.HandleFunc("/healthy", hh.IsReady)
	logger.Infof("Starting /healthz and /healthy endpoint on 0.0.0.0:%d", healthPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", healthPort))
	if err != nil {
		logger.Fatalf("Failed to start Health endpoint: %v", err)
	}
	go http.Serve(lis, m)
}
