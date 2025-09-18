package tests

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestApiTestScenarioExample demonstrates the new response body assertion capabilities
func TestApiTestScenarioExample(t *testing.T) {
	t.Parallel()

	// Example 1: Testing exact response body match
	exactMatchTest := ApiTestScenario{
		Name:           "exact match example",
		Method:         http.MethodPost,
		Url:            "/api/example",
		Body:           strings.NewReader(`{"name":"test"}`),
		ExpectedStatus: 200,
		ExpectedResponseBody: map[string]interface{}{
			"data": map[string]interface{}{
				"message": "success",
			},
		},
	}
	_ = exactMatchTest // Example usage

	// Example 2: Testing partial response body content
	partialMatchTest := ApiTestScenario{
		Name:           "partial match example",
		Method:         http.MethodGet,
		Url:            "/api/users",
		ExpectedStatus: 200,
		ExpectedResponseBodyContains: map[string]interface{}{
			"data": map[string]interface{}{
				"total": 10,
			},
		},
	}
	_ = partialMatchTest // Example usage

	// Example 3: Using custom assertion function
	customAssertTest := ApiTestScenario{
		Name:           "custom assertion example",
		Method:         http.MethodGet,
		Url:            "/api/status",
		ExpectedStatus: 200,
		ResponseBodyAssertFunc: func(t *testing.T, responseBody map[string]interface{}) {
			// Custom assertion logic
			data, exists := responseBody["data"]
			assert.True(t, exists, "Response should contain 'data' field")
			
			if dataMap, ok := data.(map[string]interface{}); ok {
				status, exists := dataMap["status"]
				assert.True(t, exists, "Data should contain 'status' field")
				assert.Equal(t, "healthy", status, "Status should be 'healthy'")
			}
		},
	}
	_ = customAssertTest // Example usage

	// Note: These are just examples to show the API. 
	// They are not meant to be run as they reference non-existent endpoints.
}