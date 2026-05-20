package tests

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/paudelgaurav/gin-integration-tests/domain/models"
	"github.com/paudelgaurav/gin-integration-tests/tests/factories"
)

func TestCreateProject(t *testing.T) {
	t.Parallel()
	s := NewSuite(t)

	category := factories.ProjectCategory.Create(t, s.DB)

	s.Client.POST("/api/v1/projects").
		JSON(gin.H{
			"name":                "Gaurav Paudel",
			"endpoint":            "https://github.com/paudelgaurav",
			"project_category_id": category.ID,
		}).
		Send().
		Status(http.StatusCreated).
		JSONPath("data.Name").Equals("Gaurav Paudel").
		JSONPath("data.ID").NotEmpty()

	s.AssertCount(&models.Project{}, 1)
}

func TestCreateProject_InvalidCategory(t *testing.T) {
	t.Parallel()
	s := NewSuite(t)

	s.Client.POST("/api/v1/projects").
		JSON(gin.H{
			"name":                "Invalid Project",
			"endpoint":            "https://example.com",
			"project_category_id": 12222,
		}).
		Send().
		Status(http.StatusBadRequest)

	s.AssertCount(&models.Project{}, 0)
}

func TestPingProjects(t *testing.T) {
	t.Parallel()
	s := NewSuite(t)

	// Seed two projects; the handler kicks off outbound pings asynchronously
	// (go h.service.PingProjects), so we register stubs to keep them silent.
	category := factories.ProjectCategory.Create(t, s.DB)
	factories.Project.CreateN(t, s.DB, 2, func(p *models.Project) {
		p.ProjectCategoryID = category.ID
	})

	s.HTTP.OnGet("=~^https://example.com/").Reply(200).Empty()

	s.Client.GET("/api/v1/projects/ping").
		Send().
		Status(http.StatusOK).
		JSONPath("data").Len(2)
}
