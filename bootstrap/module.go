package bootstrap

import (
	"github.com/paudelgaurav/gin-integration-tests/domain"
	"github.com/paudelgaurav/gin-integration-tests/pkg"
	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	pkg.Module,
	domain.Module,
)
