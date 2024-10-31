package project

import (
	"github.com/paudelgaurav/gin-integration-tests/pkg/framework"
	"github.com/paudelgaurav/gin-integration-tests/pkg/infrastructure"
)

type ProjectRoute struct {
	logger  framework.Logger
	router  *infrastructure.Router
	handler *ProjectHandler
}

func NewProjectRoute(
	logger framework.Logger,
	router *infrastructure.Router,
	handler *ProjectHandler,

) {
	r := ProjectRoute{
		logger:  logger,
		router:  router,
		handler: handler,
	}

	r.Setup()
}

func (r *ProjectRoute) Setup() {
	projectRoute := r.router.V1RouterGroup.Group("/projects")

	projectRoute.POST("", r.handler.CreateProject)
	projectRoute.GET("/ping", r.handler.PingProjects)
}
