package domain

import (
	"github.com/paudelgaurav/gin-integration-tests/domain/project"
	"go.uber.org/fx"
)

var Module = fx.Options(
	project.Module,
)
