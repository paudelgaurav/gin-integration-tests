// Package gintest provides a Django-style integration testing harness for
// Gin + GORM applications wired with uber-go/fx.
//
// A typical test:
//
//	func TestCreateProject(t *testing.T) {
//	    s := gintest.New(t,
//	        gintest.WithModules(bootstrap.CommonModules),
//	        gintest.WithMigrations(func(db *gorm.DB) error {
//	            return db.AutoMigrate(&models.Project{}, &models.ProjectCategory{})
//	        }),
//	        gintest.WithDBDecorator(func(tx *gorm.DB) *infrastructure.Database {
//	            return &infrastructure.Database{DB: tx}
//	        }),
//	    )
//
//	    category := factories.ProjectCategory(s.DB).Create()
//
//	    s.Client.POST("/api/v1/projects").
//	        JSON(gin.H{"name": "Foo", "project_category_id": category.ID}).
//	        Send().
//	        Status(http.StatusCreated).
//	        JSONPath("data.name").Equals("Foo")
//
//	    s.AssertCount(&models.Project{}, 1)
//	}
package gintest

import (
	"context"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"gorm.io/gorm"
)

// setupMu serializes fx app construction across parallel suites. Many real
// apps initialise package-globals during fx construction (viper, zap, jwk
// caches, etc.) that are not concurrency-safe. After Start() returns, tests
// run independently, so this only adds tens of milliseconds of contention
// per parallel test — well worth the loss of "first-run flake".
var setupMu sync.Mutex

// Suite is the per-test harness. It exposes the running Gin engine, the
// transactional DB, an HTTP client, and an outbound HTTP mock. All resources
// are cleaned up automatically via t.Cleanup.
type Suite struct {
	T      *testing.T
	Engine *gin.Engine
	DB     *gorm.DB
	Client *Client
	HTTP   *HTTPMock

	cfg *config
	app *fxtest.App
	tx  *gorm.DB
}

// New constructs a Suite. It opens a base DB, runs migrations, begins a
// transaction, builds the user's fx graph (with the transactional DB injected),
// starts the app, and registers cleanup hooks.
func New(t *testing.T, opts ...Option) *Suite {
	t.Helper()

	cfg := &config{
		openDB: defaultDBOpener,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	baseDB, err := cfg.openDB()
	if err != nil {
		t.Fatalf("gintest: open db: %v", err)
	}

	if cfg.migrate != nil {
		if err := cfg.migrate(baseDB); err != nil {
			t.Fatalf("gintest: migrate: %v", err)
		}
	}

	tx := baseDB.Begin()
	if tx.Error != nil {
		t.Fatalf("gintest: begin tx: %v", tx.Error)
	}

	s := &Suite{
		T:   t,
		DB:  tx,
		tx:  tx,
		cfg: cfg,
	}
	s.HTTP = newHTTPMock(t)
	s.Client = newClient(s)

	// Build fx graph with the user's modules and our decorator that swaps
	// in the transactional DB. fx.Populate extracts the *gin.Engine.
	fxOpts := []fx.Option{}
	fxOpts = append(fxOpts, cfg.modules...)
	if cfg.dbDecorator != nil {
		fxOpts = append(fxOpts, cfg.dbDecorator(tx))
	}
	fxOpts = append(fxOpts, cfg.extraOptions...)
	if cfg.engineExtractor != nil {
		fxOpts = append(fxOpts, cfg.engineExtractor(&s.Engine))
	} else {
		fxOpts = append(fxOpts, fx.Populate(&s.Engine))
	}

	if cfg.silenceFxLogs {
		fxOpts = append(fxOpts, fx.NopLogger)
	}

	setupMu.Lock()
	s.app = fxtest.New(t, fxOpts...)
	s.app.RequireStart()
	setupMu.Unlock()

	t.Cleanup(func() {
		s.app.RequireStop()
		_ = tx.Rollback()
		if sqlDB, err := baseDB.DB(); err == nil {
			_ = sqlDB.Close()
		}
		s.HTTP.deactivate()
	})

	return s
}

// Context returns a fresh background context. A convenience for tests.
func (s *Suite) Context() context.Context {
	return context.Background()
}
