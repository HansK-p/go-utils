package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func LoadUrlJsonWithHttpRequest(logger *log.Entry, client *http.Client, req *http.Request, data interface{}) error {
	logger = logger.WithFields(log.Fields{"Function": "LoadUrlJsonWithHttpRequest"})
	logger.Debugf("Executing http request")
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("when doing the request: %w", err)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("when reading response body: %w", err)
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("the endpoint '%v' returned http status code '%d' and body '%s'", req.URL, res.StatusCode, string(body))
	}
	logger.Debugf("Body to Json decode: %s", string(body))
	err = json.Unmarshal(body, data)
	if err != nil {
		return fmt.Errorf("unmarshal the response body: %w", err)
	}
	logger.Debugf("Json decoded body: %#v", data)
	return nil
}

func LoadUrlJsonWithHttpClient(logger *log.Entry, client *http.Client, url string, data interface{}) error {
	logger = logger.WithFields(log.Fields{"Function": "LoadUrlJsonWithHttpClient", "URL": url})
	logger.Debugf("Creating http request")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("when creating get request")
	}
	req.Header.Add("Accept", "application/json")
	return LoadUrlJsonWithHttpRequest(logger, client, req, data)
}

func LoadUrlJson(logger *log.Entry, url string, data interface{}) error {
	logger = logger.WithFields(log.Fields{"Function": "LoadUrlJson", "URL": url})
	logger.Debugf("Creating http client")
	client := &http.Client{}
	return LoadUrlJsonWithHttpClient(logger, client, url, data)
}
