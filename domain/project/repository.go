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

func (r *ProjectRepository) IsProjectCateogoryValid(id uint) (exists bool, err error) {
	return exists, r.Model(&models.ProjectCategory{}).Select("count(*) > 0 ").Where("id = ?", id).Find(&exists).Error
}
