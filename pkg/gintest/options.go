package gintest

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

// Option configures a Suite.
type Option func(*config)

type config struct {
	modules         []fx.Option
	extraOptions    []fx.Option
	openDB          func() (*gorm.DB, error)
	migrate         func(*gorm.DB) error
	dbDecorator     func(tx *gorm.DB) fx.Option
	engineExtractor func(target **gin.Engine) fx.Option
	authProvider    func(any) (string, string)
	silenceFxLogs   bool
}

// WithModules supplies the application's fx modules. The library will compose
// these together with its own decorators to build the test app.
func WithModules(opts ...fx.Option) Option {
	return func(c *config) {
		c.modules = append(c.modules, opts...)
	}
}

// WithFxOptions appends arbitrary fx options to the test graph. Useful for
// supplying extra mocks or overriding providers.
func WithFxOptions(opts ...fx.Option) Option {
	return func(c *config) {
		c.extraOptions = append(c.extraOptions, opts...)
	}
}

// WithDBOpener overrides how the underlying base DB is created. By default,
// the suite opens an isolated in-memory SQLite database. Use this to point
// tests at a Postgres or MySQL instance (e.g. via testcontainers).
func WithDBOpener(open func() (*gorm.DB, error)) Option {
	return func(c *config) {
		c.openDB = open
	}
}

// WithMigrations runs migrations against the base DB after open and before
// the per-test transaction begins. Typically: db.AutoMigrate(&Model{}, ...).
func WithMigrations(fn func(*gorm.DB) error) Option {
	return func(c *config) {
		c.migrate = fn
	}
}

// WithDBDecorator tells the library how to wrap the transactional *gorm.DB
// into the application's DB type so fx can inject it. For example, if your
// app uses *infrastructure.Database{*gorm.DB}, pass a wrapper that constructs
// one from the supplied tx.
func WithDBDecorator[T any](wrap func(tx *gorm.DB) T) Option {
	return func(c *config) {
		c.dbDecorator = func(tx *gorm.DB) fx.Option {
			return fx.Decorate(func(_ T) T { return wrap(tx) })
		}
	}
}

// WithEngineFrom tells the library how to find the *gin.Engine inside the
// application's DI graph. Many apps wrap the engine in a router struct
// (e.g. *infrastructure.Router{Engine *gin.Engine}); without this option,
// the library tries to resolve *gin.Engine directly from fx.
func WithEngineFrom[T any](extract func(T) *gin.Engine) Option {
	return func(c *config) {
		c.engineExtractor = func(target **gin.Engine) fx.Option {
			return fx.Invoke(func(dep T) { *target = extract(dep) })
		}
	}
}

// WithAuthProvider registers a callback the auth helpers use to convert a
// user-like value into a header (name, value) pair. Returning an empty name
// disables injection. Used by Client.AsUser().
func WithAuthProvider(fn func(any) (string, string)) Option {
	return func(c *config) {
		c.authProvider = fn
	}
}

// WithSilentFxLogs suppresses fx lifecycle logging during tests.
func WithSilentFxLogs() Option {
	return func(c *config) {
		c.silenceFxLogs = true
	}
}
