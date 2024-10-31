package cmd

import (
	"context"

	"github.com/paudelgaurav/gin-integration-tests/bootstrap"
	"github.com/paudelgaurav/gin-integration-tests/pkg/framework"
	"github.com/paudelgaurav/gin-integration-tests/pkg/infrastructure"
	"go.uber.org/fx"
)

var Modules = fx.Options(
	bootstrap.CommonModules,
)

func Execute() {
	app := fx.New(
		fx.Options(
			Modules,
		),
		fx.Invoke(startWebServer),
	)

	app.Run()
}

func startWebServer(lifecycle fx.Lifecycle, router *infrastructure.Router, logger framework.Logger) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go router.RunServer()
				return nil
			},
			OnStop: func(context context.Context) error {
				return nil
			},
		},
	)
}
