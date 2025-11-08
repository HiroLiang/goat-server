package features

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"

	"github.com/cucumber/godog"
)

type APITestContext struct {
	Server       *httptest.Server
	Client       *http.Client
	BaseURL      string
	Response     *http.Response
	ResponseBody []byte
	Token        string
}

// InitializeScenario runs once before each scenario
// Set default steps of api requests
func InitializeScenario(ctx *godog.ScenarioContext) {
	apiCtx := &APITestContext{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}

	ctx.Step(`^I send a "([^"]*)" request to "([^"]*)"$`, apiCtx.iSendRequest)
	ctx.Step(`^the response should contain (.*)$`, apiCtx.theResponseShouldContain)
}

// iSendRequest Step: sends a request to the server
func (a *APITestContext) iSendRequest(method, path string) error {
	req, err := http.NewRequest(method, a.BaseURL+path, nil)
	if err != nil {
		return err
	}

	resp, err := a.Client.Do(req)
	if err != nil {
		return err
	}
	a.Response = resp
	a.ResponseBody, _ = io.ReadAll(resp.Body)
	return nil
}

// theResponseShouldContain Step: checks the response body
func (a *APITestContext) theResponseShouldContain(expectedJSON string) error {
	expectedJSON = strings.TrimSpace(expectedJSON)

	// Try to analyze the response as JSON.
	var expected, actual interface{}
	if err := json.Unmarshal([]byte(expectedJSON), &expected); err != nil {
		return fmt.Errorf("expected value is not valid JSON: %v", err)
	}
	if err := json.Unmarshal(a.ResponseBody, &actual); err != nil {
		return fmt.Errorf("response is not valid JSON: %v\nBody: %s", err, string(a.ResponseBody))
	}

	// Compare the values.
	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("expected %+v, got %+v", expected, actual)
	}

	return nil
}
