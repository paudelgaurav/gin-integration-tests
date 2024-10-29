package project

type CreateProjectRequest struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
}
