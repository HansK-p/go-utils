/*
This is not a "real" unit test, I assume, as it actually starts a real web server on a real port - and takes it down again.
*/
package utils

import (
	"log"
	"net/http"
	"testing"
)

type MockHealth struct {
}

func (mh *MockHealth) IsAlive() (bool, string) {
	return true, "I'm alive"
}
func (mh *MockHealth) IsReady() (bool, string) {
	return true, "I'm ready"
}

func validate(t *testing.T, url string) {
	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("Error calling %s: %v", url, err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Expected http response code 200, got %v for URL: %v", resp.StatusCode, url)
	}
	t.Logf("Response: %#v", resp)
}

func TestRunHTTPHealthListener(t *testing.T) {
	log.Println("Hmmm")
	mh := MockHealth{}
	phh := PodHealthHandler{PodHealthObject: &mh}
	log.Printf("PHH: %v", phh)
	logger := GetLogger()
	RunPodHTTPHealthListener(logger, &phh)
	validate(t, "http://localhost:8080/healthy")
	validate(t, "http://localhost:8080/healthz")
}
