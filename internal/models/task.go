package models

import (
	"time"
)

// Task represents a task/todo item
type Task struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	Status      string    `json:"status" gorm:"default:'pending'"`
	Priority    int       `json:"priority" gorm:"default:1"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Version     int       `json:"version" gorm:"default:1"`
}

// TaskVersion represents a version of a task
type TaskVersion struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	TaskID    uint      `json:"task_id"`
	Title     string    `json:"title"`
	Description string   `json:"description"`
	Status    string    `json:"status"`
	Priority  int       `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
	Version   int       `json:"version"`
}

// TaskFilter represents filtering options for tasks
type TaskFilter struct {
	Status     string
	Priority   *int
	DateFrom   *time.Time
	DateTo     *time.Time
	SearchText string
}