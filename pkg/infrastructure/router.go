package infrastructure

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paudelgaurav/gin-boilerplate/pkg/framework"
)

// Router -> Gin Router
type Router struct {
	Engine        *gin.Engine
	V1RouterGroup *gin.RouterGroup
}

// NewRouter : all the routes are defined here
func NewRouter(
	env framework.Env,
) *Router {

	appEnv := env.Environment
	if appEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	httpRouter := gin.Default()

	httpRouter.MaxMultipartMemory = env.MaxMultipartMemory

	httpRouter.GET("/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "healthy"})
	})

	v1RouterGroup := httpRouter.Group("/api/v1")

	return &Router{
		Engine:        httpRouter,
		V1RouterGroup: v1RouterGroup,
	}
}

func (r *Router) RunServer() {
	if err := r.Engine.Run(); err != nil {
		log.Fatal(err)
	}
}
