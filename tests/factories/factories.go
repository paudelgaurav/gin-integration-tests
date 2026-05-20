// Package factories provides reusable model builders for tests, in the
// spirit of Django's factory-boy.
package factories

import (
	"fmt"

	"github.com/paudelgaurav/gin-integration-tests/domain/models"
	"github.com/paudelgaurav/gin-integration-tests/pkg/gintest"
)

// ProjectCategory builds a unique ProjectCategory each call.
var ProjectCategory = gintest.NewFactory(func(seq int) models.ProjectCategory {
	return models.ProjectCategory{
		Name: fmt.Sprintf("Category %d", seq),
	}
})

// Project builds a Project with a placeholder endpoint. Callers should set
// ProjectCategoryID via an override or first create a ProjectCategory.
var Project = gintest.NewFactory(func(seq int) models.Project {
	return models.Project{
		Name:     fmt.Sprintf("Project %d", seq),
		Endpoint: fmt.Sprintf("https://example.com/project/%d", seq),
	}
})
