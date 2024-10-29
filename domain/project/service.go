package project

import (
	"log"
	"net/http"

	"github.com/paudelgaurav/gin-boilerplate/domain/models"
)

type ProjectService struct {
	repo *ProjectRepository
}

func NewProjectService(repo *ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) CreateProject(projectRequest CreateProjectRequest) (*models.Project, error) {

	project := &models.Project{
		Name:     projectRequest.Name,
		Endpoint: projectRequest.Endpoint,
	}

	if err := s.repo.Create(project).Error; err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) GetAllProjects() ([]models.Project, error) {
	return s.repo.GetAllProjects()
}

func (s *ProjectService) PingProjects(projects []models.Project) {
	for i := range projects {
		endpoint := projects[i].Endpoint

		resp, err := http.Get(endpoint)
		if err != nil {
			log.Print("ping failed", err)
		}
		defer resp.Body.Close()

		log.Print(resp.StatusCode)

	}

}
