package project

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(
		NewProjectRepository,
		NewProjectService,
		NewProjectHandler,
	),
	fx.Invoke(NewProjectRoute),
)
