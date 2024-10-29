package tests

import (
	"net/http"
	"strings"
	"testing"
)

func TestPing(t *testing.T) {
	t.Parallel()

	testCases := []ApiTestScenario{
		{
			Name:           "ping",
			Method:         http.MethodGet,
			Url:            "/api/v1/projects/ping",
			ExpectedStatus: 200,
		},
		{
			Name:           "create",
			Method:         http.MethodPost,
			Url:            "/api/v1/projects",
			Body:           strings.NewReader(`{"name": "Gaurav", "endpoint": "https://github.com/paudelgaurav"}`),
			ExpectedStatus: 201,
		},
	}

	for _, testCase := range testCases {
		testCase.Test(t)
	}

}
