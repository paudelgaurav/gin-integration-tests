package gintest

import (
	"fmt"
	"sync/atomic"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var sqliteCounter atomic.Uint64

// defaultDBOpener returns an isolated in-memory SQLite DB. Each call uses a
// unique shared-cache name so parallel suites don't interfere.
func defaultDBOpener() (*gorm.DB, error) {
	n := sqliteCounter.Add(1)
	dsn := fmt.Sprintf("file:gintest_%d?mode=memory&cache=shared", n)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("gintest: open sqlite: %w", err)
	}

	// Enable FK enforcement to mirror typical production setups.
	if err := db.Exec("PRAGMA foreign_keys = ON;").Error; err != nil {
		return nil, fmt.Errorf("gintest: enable fk: %w", err)
	}

	return db, nil
}
