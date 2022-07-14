package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
)

type MockData struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

const (
	mockName = "Mr & Ms Mock"
	mockAge  = 42
)

func mockFetchDataEndpoint(w http.ResponseWriter, r *http.Request, t *testing.T) {
	if accept, ok := r.Header["Accept"]; ok {
		t.Logf("Accept header in request: %v", accept)
		jsonOk := false
		for _, accepted := range accept {
			if accepted == "application/json" {
				jsonOk = true
			}
		}
		if !jsonOk {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("only json content supported"))
		}
	}
	mockData := &MockData{
		Name: mockName,
		Age:  mockAge,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mockData)
}

func TestLoadUrlJson(t *testing.T) {
	t.Logf("start mock http server")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch strings.TrimSpace(r.URL.Path) {
		case "/":
			func(w http.ResponseWriter, r *http.Request) { mockFetchDataEndpoint(w, r, t) }(w, r)
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))

	t.Logf("perform the test")
	logger := log.New().WithFields(log.Fields{})
	mockData := &MockData{}
	if err := LoadUrlJson(logger, server.URL, mockData); err != nil {
		t.Fatalf("when loading the URL: %s", err)
	}
	if mockData.Name != mockName {
		t.Errorf("Mock name should have been '%s', but was '%s'", mockName, mockData.Name)
	}
	if mockData.Age != mockAge {
		t.Errorf("Mock name should have been '%d', but was '%d'", mockAge, mockData.Age)
	}
}
