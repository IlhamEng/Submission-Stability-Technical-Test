# Stability Team Technical Test

This repository contains a simple Task Manager API built with Go and Fiber.
Your task is to improve the stability and correctness of this system.

## Setup

Install dependencies:
```bash
go mod tidy
```

Run the server:
```bash
go run main.go
```

Server will run at:
http://localhost:3000

## Available Endpoints

- `GET /tasks` : Get all tasks
- `GET /tasks/:id` : Get specific task by ID
- `POST /tasks` : Create a new task
- `PUT /tasks/:id` : Update an existing task
- `DELETE /tasks/:id` : Delete a task

---
## Project Review & Implemented Solutions

### 🐛 1. Issues Found (Bugs & Bad Practices)
After analyzing the original codebase, several critical bugs and structural issues were identified:
1. **Critical Data Race**: The `Tasks` slice was defined globally and accessed concurrently by Fiber handlers without any Mutex (`sync.Mutex`). This causes Data Race conditioning that leads to severe application crashes.
2. **Panic on Delete (Slice Out of Bounds)**: The loop iterating in `DeleteTask` did not terminate (`break` or `return`) after appending the slice when the task was found. As the slice dynamically shrank, searching subsequent iterations threw an `Out of Bounds` fatal panic.
3. **No Auto-Increment ID**: When a user creates a new task without providing an ID, it defaulted to `0`. Continual creation overrides or causes duplicate `0` IDs.
4. **Error Ignorance (String to Int Conversion)**: The code `id, _ := strconv.Atoi(idParam)` ignored the error value. If a user passed alphabetical letters as `:id`, it silently resolved to ID `0`, potentially deleting/fetching incorrect data.
5. **Inaccurate HTTP Status Codes & Responses**: E.g., deleting a non-existent task still responded with `{"message": "deleted"}`, and fetching a non-existent task returned an HTTP `200 OK` status instead of `404 Not Found`.

### 🛠️ 2. Fixes Implemented
1. **Implemented `sync.RWMutex`**: Ensured all `store` operations (Read & Write) to the `Tasks` memory slice are strictly locked and unlocked safely to avoid Race Conditions.
2. **Fixed Slice Panic**: Added an immediate `return true` loop termination logic upon successfully splicing arrays in the `DeleteTask` function.
3. **Automated Resource IDs**: Created a safe central auto-incrementer global integer `nextID` in the store state. Every new `Task` generated is sequentially assigned starting from `3`.
4. **Fixed Memory References Mutable Loophole**: Ensured the `GetAllTasks` and `GetTaskByID` endpoints return a *copy/clone* value instead of sending back direct pointer memories to avoid unintentional mutations by the caller functions.

### ✨ 3. Improvements Made
1. **Input Validation & Proper Error Handling**: Added `err` check when extracting numbers from URLs to strictly return a `400 Bad Request` if the format is invalid. We also check if the request `{"title": ""}` is completely empty.
2. **Standardized REST API Responses**: Standardized JSON errors logic and Status Codes:
   - Responds `404 Not Found` if Task ID doesn't exist during fetching, deleting, or updating.
   - Responds `201 Created` specifically upon creating data successfully.
3. **Added "Update Task" Feature**: Improved the feature richness by appending a new REST Endpoint to update an existing memory object via Method `PUT /tasks/:id`.