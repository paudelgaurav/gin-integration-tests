package tests

import (
	"io"
	"net/http"
	"testing"

	"github.com/paudelgaurav/gin-integration-tests/domain/models"
	"github.com/paudelgaurav/gin-integration-tests/domain/project"
	"github.com/paudelgaurav/gin-integration-tests/pkg/infrastructure"
	"github.com/paudelgaurav/gin-integration-tests/pkg/utils"
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
		},
		{
			Name:           "invalid create",
			Method:         http.MethodPost,
			Url:            "/api/v1/projects",
			PrepareBody:    getInvalidCreateData,
			ExpectedStatus: 400,
		},
		{
			Name:           "ping",
			Method:         http.MethodGet,
			Url:            "/api/v1/projects/ping",
			ExpectedStatus: 200,
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
