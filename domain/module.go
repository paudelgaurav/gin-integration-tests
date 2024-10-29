package domain

import (
	"github.com/paudelgaurav/gin-boilerplate/domain/project"
	"go.uber.org/fx"
)

var Module = fx.Options(
	project.Module,
)
