//go:build e2e

package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/test"
)

// This serves as an end to end test for testing user requirements
func TestEndToEnd_HttpServer(t *testing.T) {
	// Create a new context
	ctx := context.Background()

	// Define godog options
	var opts = godog.Options{
		Format:   "pretty",
		Output:   colors.Colored(os.Stdout),
		Paths:    []string{"../../../test/feature"},
		TestingT: t,
	}

	// Read E2E_HTTP_SERVER environment variable
	serverAddress, err := internal.GetEnv("E2E_HTTP_SERVER")
	if err != nil {
		panic(err)
	}

	// Create new http driver
	driver := NewHttpDriver(serverAddress)
	if err != nil {
		panic(err)
	}

	// Run godog test suite
	status := godog.TestSuite{
		Name: "End to end tests using the http driver",
		ScenarioInitializer: test.Initialize(func() *test.Client {
			// Create a new client for each scenario (this allows to keep the client simple)
			return test.NewClient(ctx, driver)
		}),
		Options: &opts,
	}.Run()

	// Check exit status
	if status != 0 {
		os.Exit(status)
	}
}

// The HttpDriver interacts with the Http server
type HttpDriver struct {
	client   *http.Client
	endpoint string
}

func NewHttpDriver(endpoint string) *HttpDriver {
	// Create client
	client := &http.Client{CheckRedirect: CheckRedirect}

	return &HttpDriver{client: client, endpoint: endpoint}
}

// Prevent the client to follow redirections
func CheckRedirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

func (d *HttpDriver) CreateRedirection(ctx context.Context, location string) (string, error) {
	// Compose target url
	url := fmt.Sprintf("%s/api/redirection", d.endpoint)

	// Define request data
	requestData := map[string]string{
		"location": location,
	}

	// Encode data
	jsonRequestData, err := json.Marshal(requestData)
	if err != nil {
		return "", err
	}

	// Create a new request
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonRequestData))
	if err != nil {
		return "", err
	}

	// Set Content-Type header
	request.Header.Set("Content-Type", "application/json")

	// Execute the request
	response, err := d.client.Do(request)
	if err != nil {
		return "", err
	}

	// Check the status code
	if response.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("status code, expected %d, got %d", http.StatusCreated, response.StatusCode)
	}

	// Read response body
	jsonResponseData, _ := ioutil.ReadAll(response.Body)

	// Parse response body into map
	responseData := map[string]string{}
	json.Unmarshal(jsonResponseData, &responseData)

	// Return short url
	return responseData["key"], nil
}
func (d *HttpDriver) DeleteRedirection(ctx context.Context, key string) error {
	// Compose target url
	url := fmt.Sprintf("%s/api/redirection/%s", d.endpoint, key)

	// Create a new request
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	// Execute the request
	response, err := d.client.Do(request)
	if err != nil {
		return err
	}

	// Check the status code
	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("status code, expected %d, got %d", http.StatusNoContent, response.StatusCode)
	}

	// Return request error
	return err
}

func (d *HttpDriver) GetRedirectionCount(ctx context.Context, key string) (int, error) {
	// Compose target url
	url := fmt.Sprintf("%s/api/redirection/%s/count", d.endpoint, key)

	// Create a new request
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	// Execute the request
	response, err := d.client.Do(request)
	if err != nil {
		return 0, err
	}

	// Check the status code
	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("status code, expected %d, got %d", http.StatusOK, response.StatusCode)
	}

	// Read response body
	jsonResponseData, _ := ioutil.ReadAll(response.Body)

	// Parse response body into map
	responseData := map[string]int{}
	json.Unmarshal(jsonResponseData, &responseData)

	// Return visit count
	return responseData["count"], nil
}

func (d *HttpDriver) GetRedirectionLocation(ctx context.Context, key string) (string, error) {
	// Compose target url
	url := fmt.Sprintf("%s/%s", d.endpoint, key)

	// Create a new request
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	// Execute the request
	response, err := d.client.Do(request)
	if err != nil {
		return "", err
	}

	// Check the status code
	if response.StatusCode != http.StatusFound {
		return "", fmt.Errorf("status code, expected %d, got %d", http.StatusFound, response.StatusCode)
	}

	// Return short url
	return response.Header.Get("Location"), nil
}

func (d *HttpDriver) GetRedirectionList(ctx context.Context) ([]string, error) {
	// Compose target url
	url := fmt.Sprintf("%s/api/redirections", d.endpoint)

	// Create a new request
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Execute the request
	response, err := d.client.Do(request)
	if err != nil {
		return nil, err
	}

	// Check the status code
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code, expected %d, got %d", http.StatusFound, response.StatusCode)
	}

	// Read response body
	jsonResponseData, _ := ioutil.ReadAll(response.Body)

	// Parse response body into map
	responseData := map[string][]map[string]string{}
	json.Unmarshal(jsonResponseData, &responseData)

	// Compute key slice
	keys := []string{}
	for _, item := range responseData["items"] {
		keys = append(keys, item["key"])
	}

	// Return visit count
	return keys, nil
}
