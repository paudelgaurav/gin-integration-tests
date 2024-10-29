package tests

import (
	"github.com/glebarez/sqlite"
	"github.com/paudelgaurav/gin-boilerplate/pkg/infrastructure"
	"gorm.io/gorm"
)

func NewTestDatabase() *infrastructure.Database {
	db, err := gorm.Open(sqlite.Open("data_test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return &infrastructure.Database{db}
}
