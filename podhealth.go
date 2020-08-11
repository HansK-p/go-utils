package utils

import (
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// PodHealthObject interface implement a metod to validate if the object is healthy
type PodHealthObject interface {
	IsAlive() (bool, string)
	IsReady() (bool, string)
}

// PodHealthHandler contains what is needed in order to validate health
type PodHealthHandler struct {
	PodHealthObject PodHealthObject
}

// IsAlive is to be used by the Pod liveness probe
func (phh *PodHealthHandler) IsAlive(w http.ResponseWriter, r *http.Request) {
	health, message := phh.PodHealthObject.IsAlive()
	if health {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(message))
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(message))
}

// IsReady is to be used by the Pod readiness probe
func (phh *PodHealthHandler) IsReady(w http.ResponseWriter, r *http.Request) {
	health, message := phh.PodHealthObject.IsReady()
	if health {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(message))
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(message))
}

//RunPodHTTPHealthListener will start a listener listening for health and liveness checks
func RunPodHTTPHealthListener(logger *log.Logger, address string, phh *PodHealthHandler) {
	m := http.NewServeMux()
	m.HandleFunc("/healthz", phh.IsAlive)
	m.HandleFunc("/healthy", phh.IsReady)
	logger.Infof("Starting /healthz and /healthy endpoints on %v", address)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		logger.Fatalf("Failed to start Health endpoint: %v", err)
	}
	go http.Serve(lis, m)
}
