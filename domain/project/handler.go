package project

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paudelgaurav/gin-boilerplate/pkg/response"
)

type ProjectHandler struct {
	service *ProjectService
}

func NewProjectHandler(service *ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {

	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	// check validity of project category id
	valid, err := h.service.IsProjectCateogoryValid(req.ProjectCategoryID)
	if err != nil {
		response.InternalServerError(c, err.Error())
	}
	if !valid {
		response.ValidationError(c, "invalid category id ")
		return
	}

	project, err := h.service.CreateProject(req)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.JSON(c, http.StatusCreated, project)

}

func (h *ProjectHandler) PingProjects(c *gin.Context) {
	projects, err := h.service.GetAllProjects()
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	go h.service.PingProjects(projects)

	response.JSON(c, http.StatusOK, projects)
}
