package tests

import (
	"context"
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

	assert.Equal(t, scenario.ExpectedStatus, res.StatusCode)
}
