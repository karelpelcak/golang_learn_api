# TaskFlow API

TaskFlow is a scalable task management API with versioning, filtering, caching, and other enterprise-grade features.

## Features

- Task creation, retrieval, update, and deletion
- Task versioning (every update creates a new version)
- Filtering by status, priority, date, and text search
- In-memory caching for improved performance
- Proper error handling and logging
- SQLite database storage

## API Endpoints

### Create a Task
`POST /tasks`

Request body:
```json
{
  "title": "Task Title",
  "description": "Task Description",
  "status": "pending",
  "priority": 1
}
```

### Get All Tasks
`GET /tasks`

Optional query parameters:
- `status` - Filter by status (pending, in_progress, completed)
- `priority` - Filter by priority (1-5)
- `search` - Search in title or description

### Get a Task
`GET /tasks/{id}`

### Update a Task
`PUT /tasks/{id}`

Request body:
```json
{
  "title": "Updated Task Title",
  "description": "Updated Task Description",
  "status": "in_progress",
  "priority": 2
}
```

### Delete a Task
`DELETE /tasks/{id}`

### Get Task Versions
`GET /tasks/{id}/versions`

## Installation

1. Clone the repository
2. Run `go mod tidy` to install dependencies
3. Run `go run cmd/taskflow/main.go` to start the server

The server will start on port 8080 by default.