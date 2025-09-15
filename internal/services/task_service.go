package services

import (
	"taskflow/internal/models"
	"taskflow/internal/repositories"
	"taskflow/pkg/cache"
	"taskflow/pkg/logging"
	"errors"
	"strconv"
	"time"
)

// TaskService handles task business logic
type TaskService struct {
	repo   *repositories.TaskRepository
	logger *logging.Logger
	cache  *cache.Cache
}

// NewTaskService creates a new task service
func NewTaskService(repo *repositories.TaskRepository, logger *logging.Logger, cache *cache.Cache) *TaskService {
	return &TaskService{
		repo:   repo,
		logger: logger,
		cache:  cache,
	}
}

// CreateTask creates a new task
func (s *TaskService) CreateTask(task *models.Task) (*models.Task, error) {
	// Validate required fields
	if task.Title == "" {
		return nil, errors.New("title is required")
	}

	// Set default values
	if task.Status == "" {
		task.Status = "pending"
	}
	
	if task.Priority == 0 {
		task.Priority = 1
	}
	
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	task.Version = 1

	// Create task in repository
	if err := s.repo.Create(task); err != nil {
		s.logger.Error("Failed to create task: " + err.Error())
		return nil, err
	}
	
	// Cache the task
	cacheKey := "task_" + strconv.Itoa(int(task.ID))
	s.cache.Set(cacheKey, task, 300) // Cache for 5 minutes
	
	s.logger.Info("Created task with ID: " + strconv.Itoa(int(task.ID)))
	return task, nil
}

// GetTask retrieves a task by ID
func (s *TaskService) GetTask(id uint) (*models.Task, error) {
	// Try to get from cache first
	cacheKey := "task_" + strconv.Itoa(int(id))
	if cachedTask, found := s.cache.Get(cacheKey); found {
		s.logger.Info("Retrieved task from cache with ID: " + strconv.Itoa(int(id)))
		return cachedTask.(*models.Task), nil
	}

	// Get from repository
	task, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get task: " + err.Error())
		return nil, err
	}
	
	// Cache the task
	s.cache.Set(cacheKey, task, 300) // Cache for 5 minutes
	
	s.logger.Info("Retrieved task with ID: " + strconv.Itoa(int(id)))
	return task, nil
}

// GetAllTasks retrieves all tasks with optional filtering
func (s *TaskService) GetAllTasks(filter *models.TaskFilter) ([]models.Task, error) {
	// Generate cache key based on filter
	cacheKey := "tasks"
	if filter != nil {
		if filter.Status != "" {
			cacheKey += "_status_" + filter.Status
		}
		if filter.Priority != nil {
			cacheKey += "_priority_" + strconv.Itoa(*filter.Priority)
		}
		if filter.SearchText != "" {
			cacheKey += "_search_" + filter.SearchText
		}
	}
	
	// Try to get from cache first
	if cachedTasks, found := s.cache.Get(cacheKey); found {
		s.logger.Info("Retrieved tasks from cache with key: " + cacheKey)
		return cachedTasks.([]models.Task), nil
	}

	// Get from repository
	tasks, err := s.repo.GetAll(filter)
	if err != nil {
		s.logger.Error("Failed to get tasks: " + err.Error())
		return nil, err
	}
	
	// Cache the tasks
	s.cache.Set(cacheKey, tasks, 120) // Cache for 2 minutes
	
	s.logger.Info("Retrieved all tasks")
	return tasks, nil
}

// UpdateTask updates a task
func (s *TaskService) UpdateTask(task *models.Task) (*models.Task, error) {
	// Validate required fields
	if task.Title == "" {
		return nil, errors.New("title is required")
	}

	task.UpdatedAt = time.Now()
	
	// Update in repository
	if err := s.repo.Update(task); err != nil {
		s.logger.Error("Failed to update task: " + err.Error())
		return nil, err
	}
	
	// Invalidate cache
	cacheKey := "task_" + strconv.Itoa(int(task.ID))
	s.cache.Delete(cacheKey)
	
	// Also invalidate the tasks cache
	s.cache.Delete("tasks")
	
	s.logger.Info("Updated task with ID: " + strconv.Itoa(int(task.ID)))
	return task, nil
}

// DeleteTask deletes a task
func (s *TaskService) DeleteTask(id uint) error {
	// Delete from repository
	if err := s.repo.Delete(id); err != nil {
		s.logger.Error("Failed to delete task: " + err.Error())
		return err
	}
	
	// Invalidate cache
	cacheKey := "task_" + strconv.Itoa(int(id))
	s.cache.Delete(cacheKey)
	
	// Also invalidate the tasks cache
	s.cache.Delete("tasks")
	
	s.logger.Info("Deleted task with ID: " + strconv.Itoa(int(id)))
	return nil
}

// GetTaskVersions retrieves all versions of a task
func (s *TaskService) GetTaskVersions(taskID uint) ([]models.TaskVersion, error) {
	versions, err := s.repo.GetVersions(taskID)
	if err != nil {
		s.logger.Error("Failed to get task versions: " + err.Error())
		return nil, err
	}
	
	s.logger.Info("Retrieved versions for task with ID: " + strconv.Itoa(int(taskID)))
	return versions, nil
}