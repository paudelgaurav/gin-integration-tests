package pkg

import (
	"github.com/paudelgaurav/gin-boilerplate/pkg/framework"
	"github.com/paudelgaurav/gin-boilerplate/pkg/infrastructure"
	"github.com/paudelgaurav/gin-boilerplate/pkg/middleware"
	"go.uber.org/fx"
)

var Module = fx.Module("pkg",
	fx.Options(
		fx.Provide(
			framework.NewEnv,
			framework.GetLogger,
		),
	),
	infrastructure.Module,
	middleware.Module,
)
