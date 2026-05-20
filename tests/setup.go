package tests

import (
	"os"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/paudelgaurav/gin-integration-tests/bootstrap"
	"github.com/paudelgaurav/gin-integration-tests/domain/models"
	"github.com/paudelgaurav/gin-integration-tests/pkg/gintest"
	"github.com/paudelgaurav/gin-integration-tests/pkg/infrastructure"
	"gorm.io/gorm"
)

// envOnce writes a .env to the test cwd once per process. The app's
// framework.NewEnv calls log.Fatal if .env is missing, so this keeps tests
// runnable without per-developer setup. The file is gitignored.
var envOnce sync.Once

func ensureEnv() {
	envOnce.Do(func() {
		if _, err := os.Stat(".env"); err == nil {
			return
		}
		_ = os.WriteFile(".env", []byte("ENVIRONMENT=test\nLOG_LEVEL=panic\nSERVER_PORT=0\nTIMEZONE=UTC\n"), 0o600)
	})
}

// NewSuite returns a gintest.Suite preconfigured for this app: it boots the
// production fx graph, migrates the schema into a fresh in-memory SQLite DB,
// and wraps every request in a transaction that rolls back on cleanup.
func NewSuite(t *testing.T, opts ...gintest.Option) *gintest.Suite {
	t.Helper()
	ensureEnv()

	defaults := []gintest.Option{
		gintest.WithModules(bootstrap.CommonModules),
		gintest.WithMigrations(func(db *gorm.DB) error {
			return db.AutoMigrate(&models.ProjectCategory{}, &models.Project{})
		}),
		gintest.WithDBDecorator(func(tx *gorm.DB) *infrastructure.Database {
			return &infrastructure.Database{DB: tx}
		}),
		gintest.WithEngineFrom(func(r *infrastructure.Router) *gin.Engine {
			return r.Engine
		}),
		gintest.WithSilentFxLogs(),
	}

	return gintest.New(t, append(defaults, opts...)...)
}
