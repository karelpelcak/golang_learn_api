package main

import (
	"taskflow/internal/database"
	"taskflow/internal/handlers"
	"taskflow/internal/repositories"
	"taskflow/internal/services"
	"taskflow/pkg/cache"
	"taskflow/pkg/logging"
	"net/http"
	"os"
)

func main() {
	// Initialize database
	database.ConnectDB()

	// Initialize packages
	logger := logging.NewLogger()
	cache := cache.NewCache()

	// Initialize repositories
	taskRepo := repositories.NewTaskRepository()

	// Initialize services
	taskService := services.NewTaskService(taskRepo, logger, cache)

	// Initialize handlers
	taskHandler := handlers.NewTaskHandler(taskService)

	// Register routes
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskHandler.GetAllTasks(w, r)
		case http.MethodPost:
			taskHandler.CreateTask(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		// Check if it's a versions request
		if len(r.URL.Path) > 8 && r.URL.Path[len(r.URL.Path)-9:] == "/versions" {
			taskHandler.GetTaskVersions(w, r)
			return
		}

		// Handle task by ID
		switch r.Method {
		case http.MethodGet:
			taskHandler.GetTask(w, r)
		case http.MethodPut:
			taskHandler.UpdateTask(w, r)
		case http.MethodDelete:
			taskHandler.DeleteTask(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Get port from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger.Info("Starting server on port " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logger.Error("Failed to start server: " + err.Error())
		os.Exit(1)
	}
}