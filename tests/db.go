package tests

import (
	"github.com/glebarez/sqlite"
	"github.com/paudelgaurav/gin-integration-tests/domain/models"
	"github.com/paudelgaurav/gin-integration-tests/pkg/infrastructure"
	"gorm.io/gorm"
)

func NewTestDatabase() *infrastructure.Database {
	db, err := gorm.Open(sqlite.Open("data_test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Enable foreign key constraints in SQLite
	db.Exec("PRAGMA foreign_keys = ON;")

	migrate(db)

	return &infrastructure.Database{DB: db}
}

func migrate(db *gorm.DB) {
	if err := db.AutoMigrate(&models.ProjectCategory{}); err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&models.Project{}); err != nil {
		panic(err)
	}

}
