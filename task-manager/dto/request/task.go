package request

import (
	"time"
)

type CreateTaskRequest struct {
	Title        string     `json:"title" binding:"required"`
	Description  string     `json:"description"`
	AssignedToID uint       `json:"assigned_to_id" binding:"required"`
	DueDate      *time.Time `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title        *string    `json:"title"`
	Description  *string    `json:"description"`
	Status       *string    `json:"status"`
	AssignedToID *uint      `json:"assigned_to_id"`
	DueDate      *time.Time `json:"due_date"`
}