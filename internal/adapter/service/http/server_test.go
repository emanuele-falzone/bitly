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
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/emanuelefalzone/bitly/internal"
	"github.com/emanuelefalzone/bitly/test/acceptance/client"
	"github.com/emanuelefalzone/bitly/test/acceptance/driver"
	"github.com/emanuelefalzone/bitly/test/acceptance/scenario"
)

/*
This serves as an end to end test for testing user requirements
*/

func TestAcceptance_HttpDriver(t *testing.T) {

	ctx := context.Background()

	var opts = godog.Options{
		Format:   "pretty",
		Output:   colors.Colored(os.Stdout),
		Paths:    []string{"../../../../test/acceptance/feature"},
		TestingT: t,
	}

	serverAddress, err := internal.GetEnv("E2E_HTTP_SERVER")
	if err != nil {
		panic(err)
	}

	driver_ := NewHttpDriver(serverAddress)
	if err != nil {
		panic(err)
	}

	status := godog.TestSuite{
		Name: "Acceptance tests using go driver and redis repository",
		ScenarioInitializer: scenario.Initialize(func() *client.Client {
			return client.NewClient(driver_, ctx)
		}),
		Options: &opts,
	}.Run()

	if status != 0 {
		os.Exit(status)
	}
}

const Timeout = time.Duration(1) * time.Second

type HttpDriver struct {
	client   *http.Client
	endpoint string
}

func NewHttpDriver(endpoint string) driver.Driver {
	// Create client
	client := &http.Client{CheckRedirect: CheckRedirect}

	return &HttpDriver{client: client, endpoint: endpoint}
}

func CheckRedirect(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

func (d *HttpDriver) CreateRedirection(ctx context.Context, location string) (string, error) {
	// Compose target url
	url := fmt.Sprintf("%s/api", d.endpoint)

	// Define request data
	reqeuestData := map[string]string{
		"location": location,
	}

	// Encode data
	jsonReqeuestData, err := json.Marshal(reqeuestData)
	if err != nil {
		return "", err
	}

	// Create a new request
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonReqeuestData))
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

	// Compute key from location header
	key := response.Header.Get("Location")[1:]

	// Return short url
	return key, nil
}
func (d *HttpDriver) DeleteRedirection(ctx context.Context, key string) error {
	// Compose target url
	url := fmt.Sprintf("%s/api/%s", d.endpoint, key)

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
	url := fmt.Sprintf("%s/api/%s/count", d.endpoint, key)

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
