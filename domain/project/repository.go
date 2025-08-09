package project

import (
	"github.com/paudelgaurav/gin-integration-tests/domain/models"
	"github.com/paudelgaurav/gin-integration-tests/pkg/infrastructure"
)

type ProjectRepository struct {
	*infrastructure.Database
}

func NewProjectRepository(db *infrastructure.Database) *ProjectRepository {
	return &ProjectRepository{db}
}

func (r *ProjectRepository) GetAllProjects() (projects []models.Project, err error) {
	return projects, r.Preload("ProjectCategory").Find(&projects).Error
}
