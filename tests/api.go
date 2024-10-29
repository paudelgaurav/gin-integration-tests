package tests

import (
	"context"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/joho/godotenv"
	"github.com/paudelgaurav/gin-boilerplate/bootstrap"
	"github.com/paudelgaurav/gin-boilerplate/pkg/infrastructure"
	"go.uber.org/fx"
)

type ApiTestScenario struct {
	Name   string
	Method string
	Url    string
	Body   io.Reader

	//expectations
	// ----------
	ExpectedStatus int
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

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(scenario.Method, scenario.Url, scenario.Body)

	app := fx.New(
		fx.Options(
			bootstrap.CommonModules,
		),
		fx.Invoke(func(router *infrastructure.Router) {
			router.Engine.ServeHTTP(recorder, req)
		}),
	)

	// Use context.WithCancel for manual control over cancellation
	startCtx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure resources are released when the test is done

	if err := app.Start(startCtx); err != nil {
		t.Fatalf("Failed to start fx app: %v", err)
	}
	defer app.Stop(startCtx)

	// Verify the response
	res := recorder.Result()
	defer res.Body.Close()

	if res.StatusCode != scenario.ExpectedStatus {
		t.Errorf("Expected status code %d, got %d", scenario.ExpectedStatus, res.StatusCode)
	}
}
