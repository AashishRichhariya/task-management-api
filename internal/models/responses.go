package models

// Standard API responses
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type PaginationMeta struct {
	Page    int  `json:"page"`
	Limit   int  `json:"limit"`
	Total   int  `json:"total"`
	Pages   int  `json:"pages"`
	HasNext bool `json:"has_next"`
	HasPrev bool `json:"has_prev"`
}

type PaginatedTasksResponse struct {
	Tasks      []Task        `json:"tasks"`
	Pagination PaginationMeta `json:"pagination"`
}

// Task-specific responses
type TaskResponse struct {
	Message string `json:"message"`
	Data    *Task  `json:"data"`
}

type TasksResponse struct {
	Message string                  `json:"message"`
	Data    *PaginatedTasksResponse `json:"data"`
}