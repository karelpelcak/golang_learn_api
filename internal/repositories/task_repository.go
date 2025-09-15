package repositories

import (
	"taskflow/internal/database"
	"taskflow/internal/models"
	"gorm.io/gorm"
)

// TaskRepository handles database operations for tasks
type TaskRepository struct {
	db *gorm.DB
}

// NewTaskRepository creates a new task repository
func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		db: database.DB,
	}
}

// Create creates a new task
func (r *TaskRepository) Create(task *models.Task) error {
	// Create the task
	if err := r.db.Create(task).Error; err != nil {
		return err
	}

	// Create the first version
	version := models.TaskVersion{
		TaskID:      task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		CreatedAt:   task.CreatedAt,
		Version:     task.Version,
	}
	
	return r.db.Create(&version).Error
}

// GetByID retrieves a task by its ID
func (r *TaskRepository) GetByID(id uint) (*models.Task, error) {
	var task models.Task
	err := r.db.First(&task, id).Error
	return &task, err
}

// GetAll retrieves all tasks with optional filtering
func (r *TaskRepository) GetAll(filter *models.TaskFilter) ([]models.Task, error) {
	var tasks []models.Task
	query := r.db.Model(&models.Task{})

	// Apply filters
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	
	if filter.Priority != nil {
		query = query.Where("priority = ?", *filter.Priority)
	}
	
	if filter.DateFrom != nil {
		query = query.Where("created_at >= ?", *filter.DateFrom)
	}
	
	if filter.DateTo != nil {
		query = query.Where("created_at <= ?", *filter.DateTo)
	}
	
	if filter.SearchText != "" {
		searchTerm := "%" + filter.SearchText + "%"
		query = query.Where("title LIKE ? OR description LIKE ?", searchTerm, searchTerm)
	}

	err := query.Find(&tasks).Error
	return tasks, err
}

// Update updates a task and creates a new version
func (r *TaskRepository) Update(task *models.Task) error {
	// Get the current version of the task
	var currentTask models.Task
	if err := r.db.First(&currentTask, task.ID).Error; err != nil {
		return err
	}

	// Increment version
	task.Version = currentTask.Version + 1
	task.UpdatedAt = currentTask.UpdatedAt

	// Update the task
	if err := r.db.Save(task).Error; err != nil {
		return err
	}

	// Create a new version
	version := models.TaskVersion{
		TaskID:      task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		CreatedAt:   task.UpdatedAt, // Use updated_at as version timestamp
		Version:     task.Version,
	}
	
	return r.db.Create(&version).Error
}

// Delete deletes a task
func (r *TaskRepository) Delete(id uint) error {
	return r.db.Delete(&models.Task{}, id).Error
}

// GetVersions retrieves all versions of a task
func (r *TaskRepository) GetVersions(taskID uint) ([]models.TaskVersion, error) {
	var versions []models.TaskVersion
	err := r.db.Where("task_id = ?", taskID).Order("version asc").Find(&versions).Error
	return versions, err
}