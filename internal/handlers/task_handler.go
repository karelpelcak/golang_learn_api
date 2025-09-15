package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"taskflow/internal/models"
	"taskflow/internal/services"
	http_utils "taskflow/pkg/http"
)

type TaskHandler struct {
	service *services.TaskService
}

func NewTaskHandler(service *services.TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

// CreateTask handles POST /tasks
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http_utils.SendError(w, "Invalid request body", 400)
		return
	}

	createdTask, err := h.service.CreateTask(&task)
	if err != nil {
		http_utils.SendError(w, err.Error(), 400)
		return
	}

	http_utils.SendSuccess(w, createdTask)
}

// GetTask handles GET /tasks/{id}
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http_utils.SendError(w, "Invalid task ID", 400)
		return
	}

	task, err := h.service.GetTask(uint(id))
	if err != nil {
		http_utils.SendError(w, "Task not found", 404)
		return
	}

	http_utils.SendSuccess(w, task)
}

// GetAllTasks handles GET /tasks
func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	var filter models.TaskFilter

	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = status
	}

	if priorityStr := r.URL.Query().Get("priority"); priorityStr != "" {
		if priority, err := strconv.Atoi(priorityStr); err == nil {
			filter.Priority = &priority
		}
	}

	if searchText := r.URL.Query().Get("search"); searchText != "" {
		filter.SearchText = searchText
	}

	tasks, err := h.service.GetAllTasks(&filter)
	if err != nil {
		http_utils.SendError(w, "Failed to retrieve tasks", 500)
		return
	}

	http_utils.SendSuccess(w, tasks)
}

// UpdateTask handles PUT /tasks/{id}
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http_utils.SendError(w, "Invalid task ID", 400)
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http_utils.SendError(w, "Invalid request body", 400)
		return
	}

	task.ID = uint(id)

	updatedTask, err := h.service.UpdateTask(&task)
	if err != nil {
		http_utils.SendError(w, err.Error(), 400)
		return
	}

	http_utils.SendSuccess(w, updatedTask)
}

// DeleteTask handles DELETE /tasks/{id}
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http_utils.SendError(w, "Invalid task ID", 400)
		return
	}

	if err := h.service.DeleteTask(uint(id)); err != nil {
		http_utils.SendError(w, "Failed to delete task", 500)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Task deleted successfully",
	}
	http_utils.SendJSON(w, response, http.StatusOK)
}

// GetTaskVersions handles GET /tasks/{id}/versions
func (h *TaskHandler) GetTaskVersions(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(strings.TrimSuffix(r.URL.Path, "/versions"), "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http_utils.SendError(w, "Invalid task ID", 400)
		return
	}

	versions, err := h.service.GetTaskVersions(uint(id))
	if err != nil {
		http_utils.SendError(w, "Failed to retrieve task versions", 500)
		return
	}

	http_utils.SendSuccess(w, versions)
}
