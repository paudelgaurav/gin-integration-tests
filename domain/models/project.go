package models

import (
	_ "ariga.io/atlas-provider-gorm/gormschema"
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Name     string `gorm:"size:255"`
	Endpoint string `gorm:"size:255"`
}
