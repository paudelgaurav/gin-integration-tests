package infrastructure

import (
	"github.com/glebarez/sqlite"
	"github.com/paudelgaurav/gin-boilerplate/pkg/framework"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func NewDatabase(logger framework.Logger) *Database {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{Logger: logger.GetGormLogger()})
	if err != nil {
		logger.Panic(err)
	}

	return &Database{DB: db}
}
