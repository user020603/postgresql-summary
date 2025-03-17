package response

import (
	"task-manager/model"
	"time"
)

type TaskResponse struct {
	ID           uint       `json:"id"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Status       string     `json:"status"`
	AssignedToID uint       `json:"assigned_to_id"`
	AssignedTo   string     `json:"assigned_to,omitempty"`
	DueDate      *time.Time `json:"due_date"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type TaskListResponse struct {
	Tasks []TaskResponse `json:"tasks"`
	Total int64          `json:"total"`
}

// Map from Task model to TaskResponse
func MapTaskToResponse(task model.Task, includeUser bool) TaskResponse {
	response := TaskResponse{
		ID:           task.ID,
		Title:        task.Title,
		Description:  task.Description,
		Status:       string(task.Status),
		AssignedToID: task.AssignedToID,
		DueDate:      task.DueDate,
		CreatedAt:    task.CreatedAt,
		UpdatedAt:    task.UpdatedAt,
	}

	if includeUser && task.AssignedTo.ID != 0 {
		response.AssignedTo = task.AssignedTo.Username
	}

	return response
}