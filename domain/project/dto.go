package project

type CreateProjectRequest struct {
	Name     string `json:"name" binding:"required"`
	Endpoint string `json:"endpoint" binding:"required"`
}
