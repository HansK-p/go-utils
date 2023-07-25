package utils

import (
	"net/http"
	"testing"

	log "github.com/sirupsen/logrus"
)

type MockHealth struct{}

func (mh *MockHealth) IsAlive() (bool, string) { return true, "I'm alive" }
func (mh *MockHealth) IsReady() (bool, string) { return true, "I'm ready" }

type MockUnHealth struct{}

func (muh *MockUnHealth) IsAlive() (bool, string) { return false, "I'm not alive" }
func (muh *MockUnHealth) IsReady() (bool, string) { return false, "I'm not ready" }

func validate(t *testing.T, url string, expect int) {
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Error calling %s: %v", url, err)
	}
	if resp.StatusCode != expect {
		t.Fatalf("Expected http response code %d, got %d for URL: %v", expect, resp.StatusCode, url)
	}
	t.Logf("Response: %#v", resp)
}

func TestRunHTTPHealthListener(t *testing.T) {
	mh, muh := MockHealth{}, MockUnHealth{}
	phh := PodHealthHandler{PodHealthObject: &mh}
	logger := NewLogger().WithFields(log.Fields{})
	RunPodHTTPHealthListener(logger, "127.0.0.1:8080", &phh)
	validate(t, "http://localhost:8080/healthy", 200)
	validate(t, "http://localhost:8080/healthz", 200)

	phh.PodHealthObjects = []PodHealthObject{&mh, &mh}
	validate(t, "http://localhost:8080/healthy", 200)
	validate(t, "http://localhost:8080/healthz", 200)

	phh.PodHealthObjects = []PodHealthObject{&mh, &muh}
	validate(t, "http://localhost:8080/healthy", 500)
	validate(t, "http://localhost:8080/healthz", 500)
}
