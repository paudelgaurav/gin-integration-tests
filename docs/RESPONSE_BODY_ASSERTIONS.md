# Response Body Assertion in API Tests

This document explains how to use the enhanced `ApiTestScenario` struct to assert response body contents along with status codes.

## Overview

The `ApiTestScenario` struct now supports three types of response body assertions:

1. **Exact Response Body Matching** - `ExpectedResponseBody`
2. **Partial Content Matching** - `ExpectedResponseBodyContains`
3. **Custom Assertion Functions** - `ResponseBodyAssertFunc`

## Usage Examples

### 1. Exact Response Body Matching

Use `ExpectedResponseBody` when you want to validate the entire response structure:

```go
testCase := ApiTestScenario{
    Name:           "exact match test",
    Method:         http.MethodGet,
    Url:            "/api/status",
    ExpectedStatus: 200,
    ExpectedResponseBody: map[string]interface{}{
        "data": map[string]interface{}{
            "status": "healthy",
            "version": "1.0.0",
        },
    },
}
```

### 2. Partial Content Matching

Use `ExpectedResponseBodyContains` when you want to validate specific fields without checking the entire response:

```go
testCase := ApiTestScenario{
    Name:           "partial match test",
    Method:         http.MethodPost,
    Url:            "/api/v1/projects",
    PrepareBody:    getCreateData,
    ExpectedStatus: 201,
    ExpectedResponseBodyContains: map[string]interface{}{
        "data": map[string]interface{}{
            "Name":     "Gaurav Paudel",
            "Endpoint": "https://github.com/paudelgaurav",
        },
    },
}
```

### 3. Custom Assertion Functions

Use `ResponseBodyAssertFunc` for complex validation logic:

```go
testCase := ApiTestScenario{
    Name:           "custom assertion test",
    Method:         http.MethodGet,
    Url:            "/api/v1/projects/ping",
    ExpectedStatus: 200,
    ResponseBodyAssertFunc: func(t *testing.T, responseBody map[string]interface{}) {
        data, exists := responseBody["data"]
        assert.True(t, exists, "Response should contain 'data' field")
        
        if dataArray, ok := data.([]interface{}); ok {
            assert.GreaterOrEqual(t, len(dataArray), 1, "Should have at least one project")
            
            if len(dataArray) > 0 {
                project := dataArray[0].(map[string]interface{})
                assert.NotEmpty(t, project["ID"], "Project should have ID")
                assert.NotEmpty(t, project["Name"], "Project should have Name")
            }
        }
    },
}
```

### 4. Combining Multiple Assertion Types

You can combine different assertion types in a single test:

```go
testCase := ApiTestScenario{
    Name:           "combined assertions",
    Method:         http.MethodPost,
    Url:            "/api/v1/projects",
    PrepareBody:    getCreateData,
    ExpectedStatus: 201,
    ExpectedResponseBodyContains: map[string]interface{}{
        "data": map[string]interface{}{
            "Name": "Test Project",
        },
    },
    ResponseBodyAssertFunc: func(t *testing.T, responseBody map[string]interface{}) {
        data := responseBody["data"].(map[string]interface{})
        id := data["ID"].(float64)
        assert.Greater(t, id, float64(0), "Project ID should be greater than 0")
    },
}
```

## Response Structure Support

The framework supports the following response structures:

### Success Responses
```json
{
  "data": {
    "ID": 1,
    "Name": "Project Name",
    "Endpoint": "https://example.com"
  }
}
```

### Error Responses
```json
{
  "error": "Error message"
}
```

### Array Responses
```json
{
  "data": [
    {
      "ID": 1,
      "Name": "Project 1"
    },
    {
      "ID": 2,
      "Name": "Project 2"
    }
  ]
}
```

## Backward Compatibility

All existing tests continue to work without modification. Response body assertions are optional - if none are specified, only the status code is validated.

## Running Tests

Tests should be run from the repository root directory:

```bash
# Run all tests
go test -v github.com/paudelgaurav/gin-integration-tests/tests

# Run specific test
go test -v github.com/paudelgaurav/gin-integration-tests/tests -run TestPing
```

Note: Make sure the `.env` file exists in your project root before running tests.