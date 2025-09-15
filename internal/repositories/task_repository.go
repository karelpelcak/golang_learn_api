package repositories

import (
	"taskflow/internal/database"
	"taskflow/internal/models"

	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		db: database.DB,
	}
}

func (r *TaskRepository) Create(task *models.Task) error {
	if err := r.db.Create(task).Error; err != nil {
		return err
	}

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

func (r *TaskRepository) GetByID(id uint) (*models.Task, error) {
	var task models.Task
	err := r.db.First(&task, id).Error
	return &task, err
}

func (r *TaskRepository) GetAll(filter *models.TaskFilter) ([]models.Task, error) {
	var tasks []models.Task
	query := r.db.Model(&models.Task{})

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

func (r *TaskRepository) Update(task *models.Task) error {
	var currentTask models.Task
	if err := r.db.First(&currentTask, task.ID).Error; err != nil {
		return err
	}

	task.Version = currentTask.Version + 1
	task.UpdatedAt = currentTask.UpdatedAt

	if err := r.db.Save(task).Error; err != nil {
		return err
	}

	version := models.TaskVersion{
		TaskID:      task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		CreatedAt:   task.UpdatedAt,
		Version:     task.Version,
	}

	return r.db.Create(&version).Error
}

func (r *TaskRepository) Delete(id uint) error {
	return r.db.Delete(&models.Task{}, id).Error
}

func (r *TaskRepository) GetVersions(taskID uint) ([]models.TaskVersion, error) {
	var versions []models.TaskVersion
	err := r.db.Where("task_id = ?", taskID).Order("version asc").Find(&versions).Error
	return versions, err
}
