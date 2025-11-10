package model

type ProjectDataRequest struct {
	StartDate  string   `json:"start_date"`
	EndDate    string   `json:"end_date"`
	ProjectIDs []string `json:"project_ids"`
}

func NewProjectDataRequest() *ProjectDataRequest {
	return &ProjectDataRequest{}
}
