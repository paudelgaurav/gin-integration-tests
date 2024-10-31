package models

import (
	"gorm.io/gorm"
)

type ProjectCategory struct {
	gorm.Model
	Name string `gorm:"size:255"`
}

type Project struct {
	gorm.Model
	Name              string `gorm:"size:255"`
	Endpoint          string `gorm:"size:255"`
	ProjectCategoryID uint
	ProjectCategory   ProjectCategory `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:ProjectCategoryID;references:ID"`
}
