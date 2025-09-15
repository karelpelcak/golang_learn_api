package handlers

import (
	"taskflow/internal/models"
	"taskflow/internal/services"
	http_utils "taskflow/pkg/http"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// TaskHandler handles HTTP requests for tasks
type TaskHandler struct {
	service *services.TaskService
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(service *services.TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

// CreateTask handles POST /tasks
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http_utils.SendError(w, "Invalid request body", 400)
		return
	}

	// Create task
	createdTask, err := h.service.CreateTask(&task)
	if err != nil {
		http_utils.SendError(w, err.Error(), 400)
		return
	}

	// Send response
	http_utils.SendSuccess(w, createdTask)
}

// GetTask handles GET /tasks/{id}
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http_utils.SendError(w, "Invalid task ID", 400)
		return
	}

	// Get task
	task, err := h.service.GetTask(uint(id))
	if err != nil {
		http_utils.SendError(w, "Task not found", 404)
		return
	}

	// Send response
	http_utils.SendSuccess(w, task)
}

// GetAllTasks handles GET /tasks
func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for filtering
	var filter models.TaskFilter
	
	// Status filter
	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = status
	}
	
	// Priority filter
	if priorityStr := r.URL.Query().Get("priority"); priorityStr != "" {
		if priority, err := strconv.Atoi(priorityStr); err == nil {
			filter.Priority = &priority
		}
	}
	
	// Date filters
	// TODO: Implement date filtering
	
	// Search text filter
	if searchText := r.URL.Query().Get("search"); searchText != "" {
		filter.SearchText = searchText
	}

	// Get all tasks
	tasks, err := h.service.GetAllTasks(&filter)
	if err != nil {
		http_utils.SendError(w, "Failed to retrieve tasks", 500)
		return
	}

	// Send response
	http_utils.SendSuccess(w, tasks)
}

// UpdateTask handles PUT /tasks/{id}
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http_utils.SendError(w, "Invalid task ID", 400)
		return
	}

	// Parse request body
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http_utils.SendError(w, "Invalid request body", 400)
		return
	}
	
	// Set the ID
	task.ID = uint(id)

	// Update task
	updatedTask, err := h.service.UpdateTask(&task)
	if err != nil {
		http_utils.SendError(w, err.Error(), 400)
		return
	}

	// Send response
	http_utils.SendSuccess(w, updatedTask)
}

// DeleteTask handles DELETE /tasks/{id}
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http_utils.SendError(w, "Invalid task ID", 400)
		return
	}

	// Delete task
	if err := h.service.DeleteTask(uint(id)); err != nil {
		http_utils.SendError(w, "Failed to delete task", 500)
		return
	}

	// Send response
	response := map[string]interface{}{
		"success": true,
		"message": "Task deleted successfully",
	}
	http_utils.SendJSON(w, response, http.StatusOK)
}

// GetTaskVersions handles GET /tasks/{id}/versions
func (h *TaskHandler) GetTaskVersions(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL
	idStr := strings.TrimPrefix(strings.TrimSuffix(r.URL.Path, "/versions"), "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http_utils.SendError(w, "Invalid task ID", 400)
		return
	}

	// Get task versions
	versions, err := h.service.GetTaskVersions(uint(id))
	if err != nil {
		http_utils.SendError(w, "Failed to retrieve task versions", 500)
		return
	}

	// Send response
	http_utils.SendSuccess(w, versions)
}