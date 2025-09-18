package tests

import (
	"io"
	"net/http"
	"testing"

	"github.com/paudelgaurav/gin-integration-tests/domain/models"
	"github.com/paudelgaurav/gin-integration-tests/domain/project"
	"github.com/paudelgaurav/gin-integration-tests/pkg/infrastructure"
	"github.com/paudelgaurav/gin-integration-tests/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	t.Parallel()

	testCases := []ApiTestScenario{
		{
			Name:           "create",
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
			ResponseBodyAssertFunc: func(t *testing.T, responseBody map[string]interface{}) {
				// Verify the data wrapper exists
				data, exists := responseBody["data"]
				assert.True(t, exists, "Response should contain 'data' field")
				
				if dataMap, ok := data.(map[string]interface{}); ok {
					// Verify project ID exists and is greater than 0
					if id, exists := dataMap["ID"]; exists {
						if idFloat, ok := id.(float64); ok {
							assert.Greater(t, idFloat, float64(0), "Project ID should be greater than 0")
						}
					}
					
					// Verify timestamps exist and are not zero
					assert.NotEmpty(t, dataMap["CreatedAt"], "CreatedAt should not be empty")
					assert.NotEmpty(t, dataMap["UpdatedAt"], "UpdatedAt should not be empty")
					
					// Verify ProjectCategoryID is set
					if catID, exists := dataMap["ProjectCategoryID"]; exists {
						if catIDFloat, ok := catID.(float64); ok {
							assert.Greater(t, catIDFloat, float64(0), "ProjectCategoryID should be greater than 0")
						}
					}
				}
			},
		},
		{
			Name:           "invalid create",
			Method:         http.MethodPost,
			Url:            "/api/v1/projects",
			PrepareBody:    getInvalidCreateData,
			ExpectedStatus: 400,
			ExpectedResponseBodyContains: map[string]interface{}{
				"error": "project_category_id: invalid project category ID.",
			},
		},
		{
			Name:           "ping",
			Method:         http.MethodGet,
			Url:            "/api/v1/projects/ping",
			ExpectedStatus: 200,
			ResponseBodyAssertFunc: func(t *testing.T, responseBody map[string]interface{}) {
				// Verify the data wrapper exists
				data, exists := responseBody["data"]
				assert.True(t, exists, "Response should contain 'data' field")
				
				// Verify data is an array
				if dataArray, ok := data.([]interface{}); ok {
					// There should be at least one project (created by previous test)
					assert.GreaterOrEqual(t, len(dataArray), 1, "Should have at least one project")
					
					// Verify structure of first project if exists
					if len(dataArray) > 0 {
						if project, ok := dataArray[0].(map[string]interface{}); ok {
							// Verify required fields exist
							assert.NotEmpty(t, project["ID"], "Project should have ID")
							assert.NotEmpty(t, project["Name"], "Project should have Name")
							assert.NotEmpty(t, project["Endpoint"], "Project should have Endpoint")
							assert.NotEmpty(t, project["ProjectCategoryID"], "Project should have ProjectCategoryID")
							assert.NotEmpty(t, project["CreatedAt"], "Project should have CreatedAt")
							assert.NotEmpty(t, project["UpdatedAt"], "Project should have UpdatedAt")
						}
					}
				} else {
					t.Errorf("Expected data to be an array, got %T", data)
				}
			},
		},
	}

	for _, testCase := range testCases {
		testCase.Test(t)
	}

}

func getCreateData(db *infrastructure.Database) io.Reader {

	projectCategory := models.ProjectCategory{
		Name: "Unit testing",
	}

	if err := db.Create(&projectCategory).Error; err != nil {
		panic(err)
	}

	body := project.CreateProjectRequest{
		Name:              "Gaurav Paudel",
		Endpoint:          "https://github.com/paudelgaurav",
		ProjectCategoryID: projectCategory.ID,
	}

	return utils.StructToReader(&body)

}

func getInvalidCreateData(db *infrastructure.Database) io.Reader {
	body := project.CreateProjectRequest{
		Name:              "Invalid Project",
		Endpoint:          "https://github.com/paudelgaurav",
		ProjectCategoryID: 12222,
	}

	return utils.StructToReader(&body)

}
