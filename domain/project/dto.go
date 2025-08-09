package project

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/paudelgaurav/gin-integration-tests/domain/models"
	"gorm.io/gorm"
)

type CreateProjectRequest struct {
	Name              string `json:"name" binding:"required"`
	Endpoint          string `json:"endpoint" binding:"required"`
	ProjectCategoryID uint   `json:"project_category_id" binding:"required"`
}

func (r CreateProjectRequest) Validate(tx *gorm.DB) error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.ProjectCategoryID, validation.By(func(value interface{}) error {
			// check whether the project category id exists or not
			projectCatID := value.(uint)
			var count int64
			if err := tx.Model(&models.ProjectCategory{}).Where("id = ?", projectCatID).Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				return errors.New("invalid project category ID")
			}
			return nil
		})),
	)
}
