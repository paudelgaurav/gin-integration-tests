package project

type CreateProjectRequest struct {
	Name              string `json:"name" binding:"required"`
	Endpoint          string `json:"endpoint" binding:"required"`
	ProjectCategoryID uint   `json:"project_category_id" binding:"required"`
}
