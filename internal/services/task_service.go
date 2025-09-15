package services

import (
	"errors"
	"strconv"
	"taskflow/internal/models"
	"taskflow/internal/repositories"
	"taskflow/pkg/cache"
	"taskflow/pkg/logging"
	"time"
)

type TaskService struct {
	repo   *repositories.TaskRepository
	logger *logging.Logger
	cache  *cache.Cache
}

func NewTaskService(repo *repositories.TaskRepository, logger *logging.Logger, cache *cache.Cache) *TaskService {
	return &TaskService{
		repo:   repo,
		logger: logger,
		cache:  cache,
	}
}

func (s *TaskService) CreateTask(task *models.Task) (*models.Task, error) {
	if task.Title == "" {
		return nil, errors.New("title is required")
	}

	if task.Status == "" {
		task.Status = "pending"
	}

	if task.Priority == 0 {
		task.Priority = 1
	}

	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	task.Version = 1

	if err := s.repo.Create(task); err != nil {
		s.logger.Error("Failed to create task: " + err.Error())
		return nil, err
	}

	cacheKey := "task_" + strconv.Itoa(int(task.ID))
	s.cache.Set(cacheKey, task, 300)

	s.logger.Info("Created task with ID: " + strconv.Itoa(int(task.ID)))
	return task, nil
}

func (s *TaskService) GetTask(id uint) (*models.Task, error) {
	cacheKey := "task_" + strconv.Itoa(int(id))
	if cachedTask, found := s.cache.Get(cacheKey); found {
		s.logger.Info("Retrieved task from cache with ID: " + strconv.Itoa(int(id)))
		return cachedTask.(*models.Task), nil
	}

	task, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get task: " + err.Error())
		return nil, err
	}

	s.cache.Set(cacheKey, task, 300)

	s.logger.Info("Retrieved task with ID: " + strconv.Itoa(int(id)))
	return task, nil
}

func (s *TaskService) GetAllTasks(filter *models.TaskFilter) ([]models.Task, error) {
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

	if cachedTasks, found := s.cache.Get(cacheKey); found {
		s.logger.Info("Retrieved tasks from cache with key: " + cacheKey)
		return cachedTasks.([]models.Task), nil
	}

	tasks, err := s.repo.GetAll(filter)
	if err != nil {
		s.logger.Error("Failed to get tasks: " + err.Error())
		return nil, err
	}

	s.cache.Set(cacheKey, tasks, 120)

	s.logger.Info("Retrieved all tasks")
	return tasks, nil
}

func (s *TaskService) UpdateTask(task *models.Task) (*models.Task, error) {
	if task.Title == "" {
		return nil, errors.New("title is required")
	}

	task.UpdatedAt = time.Now()

	if err := s.repo.Update(task); err != nil {
		s.logger.Error("Failed to update task: " + err.Error())
		return nil, err
	}

	cacheKey := "task_" + strconv.Itoa(int(task.ID))
	s.cache.Delete(cacheKey)

	s.cache.Delete("tasks")

	s.logger.Info("Updated task with ID: " + strconv.Itoa(int(task.ID)))
	return task, nil
}

func (s *TaskService) DeleteTask(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		s.logger.Error("Failed to delete task: " + err.Error())
		return err
	}

	cacheKey := "task_" + strconv.Itoa(int(id))
	s.cache.Delete(cacheKey)

	s.cache.Delete("tasks")

	s.logger.Info("Deleted task with ID: " + strconv.Itoa(int(id)))
	return nil
}

func (s *TaskService) GetTaskVersions(taskID uint) ([]models.TaskVersion, error) {
	versions, err := s.repo.GetVersions(taskID)
	if err != nil {
		s.logger.Error("Failed to get task versions: " + err.Error())
		return nil, err
	}

	s.logger.Info("Retrieved versions for task with ID: " + strconv.Itoa(int(taskID)))
	return versions, nil
}
