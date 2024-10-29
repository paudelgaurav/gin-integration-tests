package project

import (
	"github.com/paudelgaurav/gin-boilerplate/domain/models"
	"github.com/paudelgaurav/gin-boilerplate/pkg/infrastructure"
)

type ProjectRepository struct {
	*infrastructure.Database
}

func NewProjectRepository(db *infrastructure.Database) *ProjectRepository {
	return &ProjectRepository{db}
}

func (r *ProjectRepository) GetAllProjects() (projects []models.Project, err error) {
	return projects, r.Find(&projects).Error
}
