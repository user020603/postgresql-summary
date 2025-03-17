package model

import (
	"gorm.io/gorm"
	"time"
)

type TaskStatus string

const (
	TaskStatusTodo       TaskStatus = "TODO"
	TaskStatusInProgress TaskStatus = "IN_PROGRESS"
	TaskStatusDone       TaskStatus = "DONE"
)

type Task struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Title        string         `gorm:"not null" json:"title"`
	Description  string         `json:"description"`
	Status       TaskStatus     `gorm:"type:varchar(20);default:'TODO'" json:"status"`
	AssignedToID uint           `json:"assigned_to_id"`
	AssignedTo   User           `gorm:"foreignKey:AssignedToID" json:"assigned_to,omitempty"`
	CreatedBy    uint           `json:"created_by"`
	DueDate      *time.Time     `json:"due_date"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}