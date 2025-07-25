package models

// Task-related requests
type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=255"`
	Description string `json:"description" binding:"max=1000"`
	Status      string `json:"status" binding:"omitempty,oneof=pending in_progress completed closed"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title" binding:"omitempty,min=1,max=255"`
	Description string `json:"description" binding:"omitempty,max=1000"`
	Status      string `json:"status" binding:"omitempty,oneof=pending in_progress completed closed"`
}

// Query parameters
type TaskQueryParams struct {
	Page      int    `form:"page"`
	Limit     int    `form:"limit"`
	Status    string `form:"status" binding:"omitempty,oneof=pending in_progress completed closed"`
	SortBy    string `form:"sort_by" binding:"omitempty,oneof=id title status created_at updated_at"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// Set defaults for query params
func (q *TaskQueryParams) SetDefaults() {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit <= 0 || q.Limit > 100 {
		q.Limit = 10
	}
	if q.SortBy == "" {
		q.SortBy = "created_at"
	}
	if q.SortOrder == "" || (q.SortOrder != "asc" && q.SortOrder != "desc"){
		q.SortOrder = "desc"
	}
}

// URL parameters
type TaskIDParam struct {
	ID int `uri:"id" binding:"required,min=1"`
}