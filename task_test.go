package main

import (
	"fmt"
	"taskflow/internal/database"
	"taskflow/internal/models"
	"taskflow/internal/repositories"
	"testing"
)

func TestTaskRepository(t *testing.T) {
	database.ConnectDB()

	repo := repositories.NewTaskRepository()

	task := &models.Task{
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      "pending",
		Priority:    1,
	}

	err := repo.Create(task)
	if err != nil {
		t.Errorf("Failed to create task: %v", err)
	}

	if task.ID == 0 {
		t.Error("Task ID should not be zero after creation")
	}

	retrievedTask, err := repo.GetByID(task.ID)
	if err != nil {
		t.Errorf("Failed to get task: %v", err)
	}

	if retrievedTask.Title != task.Title {
		t.Errorf("Expected title %s, got %s", task.Title, retrievedTask.Title)
	}

	task.Title = "Updated Task"
	task.Version = 1
	err = repo.Update(task)
	if err != nil {
		t.Errorf("Failed to update task: %v", err)
	}

	updatedTask, err := repo.GetByID(task.ID)
	if err != nil {
		t.Errorf("Failed to get updated task: %v", err)
	}

	if updatedTask.Title != "Updated Task" {
		t.Errorf("Expected title %s, got %s", "Updated Task", updatedTask.Title)
	}

	if updatedTask.Version != 2 {
		t.Errorf("Expected version %d, got %d", 2, updatedTask.Version)
	}

	versions, err := repo.GetVersions(task.ID)
	if err != nil {
		t.Errorf("Failed to get task versions: %v", err)
	}

	if len(versions) != 2 {
		t.Errorf("Expected 2 versions, got %d", len(versions))
	}

	tasks, err := repo.GetAll(&models.TaskFilter{})
	if err != nil {
		t.Errorf("Failed to get all tasks: %v", err)
	}

	if len(tasks) == 0 {
		t.Error("Expected at least one task")
	}

	err = repo.Delete(task.ID)
	if err != nil {
		t.Errorf("Failed to delete task: %v", err)
	}

	_, err = repo.GetByID(task.ID)
	if err == nil {
		t.Error("Expected error when getting deleted task")
	}

	fmt.Println("All tests passed!")
}
