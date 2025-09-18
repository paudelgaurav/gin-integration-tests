package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
	"github.com/paudelgaurav/gin-integration-tests/bootstrap"
	"github.com/paudelgaurav/gin-integration-tests/pkg/infrastructure"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

type ApiTestScenario struct {
	Name        string
	Method      string
	Url         string
	Body        io.Reader
	PrepareBody func(db *infrastructure.Database) io.Reader

	//expectations
	// ----------
	ExpectedStatus int
	// Expected response body as interface{} for exact matching
	ExpectedResponseBody interface{}
	// Expected response body contains specific values/keys
	ExpectedResponseBodyContains map[string]interface{}
	// Custom assertion function for complex response validation
	ResponseBodyAssertFunc func(t *testing.T, responseBody map[string]interface{})
}

func (scenario *ApiTestScenario) getBody(db *infrastructure.Database) io.Reader {
	if scenario.Body != nil {
		return scenario.Body
	} else if scenario.PrepareBody != nil {
		return scenario.PrepareBody(db)
	}

	return nil
}

// Test executes the api test case scenario
func (scenario *ApiTestScenario) Test(t *testing.T) {

	var name = scenario.Name
	if name == "" {
		name = fmt.Sprintf("%s:%s", scenario.Method, scenario.Url)
	}

	t.Run(name, scenario.test)

}

func (scenario *ApiTestScenario) test(t *testing.T) {

	if err := godotenv.Load(); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	db := NewTestDatabase()

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(scenario.Method, scenario.Url, scenario.getBody(db))

	app := fxtest.New(
		t,
		fx.Options(
			bootstrap.CommonModules,
			fx.Decorate(NewTestDatabase),
		),
		fx.Invoke(func(router *infrastructure.Router) {
			router.Engine.ServeHTTP(recorder, req)
		}),
	)

	startCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := app.Start(startCtx); err != nil {
		t.Fatalf("Failed to start fx app: %v", err)
	}
	defer app.Stop(startCtx)

	res := recorder.Result()
	defer res.Body.Close()

	// Assert status code
	assert.Equal(t, scenario.ExpectedStatus, res.StatusCode)

	// Read response body
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Parse response body as JSON if we have body assertions
	if scenario.ExpectedResponseBody != nil || scenario.ExpectedResponseBodyContains != nil || scenario.ResponseBodyAssertFunc != nil {
		var responseBody map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &responseBody); err != nil {
			t.Fatalf("Failed to parse response body as JSON: %v", err)
		}

		// Exact response body matching
		if scenario.ExpectedResponseBody != nil {
			assert.Equal(t, scenario.ExpectedResponseBody, responseBody)
		}

		// Partial response body matching
		if scenario.ExpectedResponseBodyContains != nil {
			for key, expectedValue := range scenario.ExpectedResponseBodyContains {
				actualValue, exists := responseBody[key]
				if !exists {
					t.Errorf("Expected response body to contain key '%s'", key)
					continue
				}

				// Handle nested object matching
				if expectedMap, ok := expectedValue.(map[string]interface{}); ok {
					if actualMap, ok := actualValue.(map[string]interface{}); ok {
						for nestedKey, nestedExpectedValue := range expectedMap {
							nestedActualValue, nestedExists := actualMap[nestedKey]
							if !nestedExists {
								t.Errorf("Expected response body to contain nested key '%s.%s'", key, nestedKey)
								continue
							}
							assert.Equal(t, nestedExpectedValue, nestedActualValue, 
								"Mismatch in nested field %s.%s", key, nestedKey)
						}
					} else {
						t.Errorf("Expected nested object for key '%s' but got %T", key, actualValue)
					}
				} else {
					assert.Equal(t, expectedValue, actualValue, "Mismatch in field %s", key)
				}
			}
		}

		// Custom assertion function
		if scenario.ResponseBodyAssertFunc != nil {
			scenario.ResponseBodyAssertFunc(t, responseBody)
		}
	}
}
