package bootstrap

import (
	"github.com/paudelgaurav/gin-boilerplate/domain"
	"github.com/paudelgaurav/gin-boilerplate/pkg"
	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	pkg.Module,
	domain.Module,
)
