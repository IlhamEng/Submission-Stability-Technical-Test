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

This section details the vulnerabilities found in the original codebase and the systematic fixes applied to ensure production-grade stability.

### 1. Issues Found (Bugs & Vulnerabilities)

| Issue | Severity | Location | Description & Impact |
| :--- | :---: | :--- | :--- |
| **Data Race Condition** | 🔴 Critical | `store/task_store.go` | The global `Tasks` slice was mutated concurrently by multiple API requests without Mutex locking. **Impact**: High risk of `panic` (application crash) under concurrent traffic. |
| **Slice Out-of-Bounds Panic** | 🔴 Critical | `store` (`DeleteTask`) | The loop iterating to delete a task did not terminate early after the deletion. As the slice length shrank, subsequent loop cycles attempted to access a non-existent index. **Impact**: `panic: slice bounds out of range`. |
| **No Auto-Increment ID** | 🟠 High | `store` / `handlers` | Creating a task without passing an `"id"` defaults to Go's zero-value `0`. Subsequent creations keep getting ID `0`, overriding or duplicating entries. **Impact**: Data corruption and duplicate records. |
| **Error Ignorance** | 🟡 Medium | `handlers` (`GetTask`, `Delete`) | The codebase ignored validation errors from `strconv.Atoi()`. If an alphabetical string was passed as an ID (e.g., `/tasks/abc`), it silently fell back to ID `0`. **Impact**: Unexpected behavior, acting on the wrong task. |
| **Improper HTTP Protocols** | 🟡 Medium | All `handlers` | Handlers responded with `200 OK` or `{"message": "deleted"}` even if the target Task ID did not exist. **Impact**: Bad API design, misleading the client frontend. |

---

### 2. Fixes Implemented

Here is how the critical issues were resolved:

#### A. Concurrency Protection (Sync Mutex)
Added `sync.RWMutex` to the store package. This ensures that only one request can write to the memory at a time, preventing Data Races.

**Before:**
```go
func AddTask(task models.Task) {
	Tasks = append(Tasks, task) // Unsafe concurrent write
}
```
**After:**
```go
var mu sync.RWMutex

func AddTask(task models.Task) models.Task {
	mu.Lock() // Exclusive lock acquired
	defer mu.Unlock()
	Tasks = append(Tasks, task) 
	return task
}
```

#### B. Fixing Slice Deletion Panic
Added a `return true` statement immediately after the element is removed from the slice, safely stopping the iteration loop.

#### C. Automated Resource ID Generator
Created a global `nextID` integer in the `store` state. Whenever `AddTask()` is called, it inherently assigns sequential IDs (`3, 4, 5...`), removing the client's burden of generating IDs manually.

#### D. Preventing Pointer Mutation (Memory Safety)
Modified the endpoints (`GetAllTasks` and `GetTaskByID`) to return a cloned/copied value rather than direct pointer memory, preventing the caller (handlers) from mutating the central `Tasks` store unintentionally.

#### E. Robust Input Validation
   - Empty title check: If a client submits `{"title": ""}`, the server correctly rejects it.
   - Numeric URL Parameter check: Safely catches alphabetical inputs and blocks them.
#### F. Standardized REST HTTP Codes
   - Now properly returns `404 Not Found` when trying to Fetch, Delete, or Update a non-existent task.
   - Now returns `201 Created` specifically upon a successful `POST` request.
   - Now returns `400 Bad Request` for malformed payloads.

---

### 3. Improvements Made

Beyond bug fixes, the following enhancements elevate the application's overall quality:


1. **New Endpoint Added (`Update Task`):**
   - Implemented an entire flow (Store Function + Handler + Fiber Route) for `PUT /tasks/:id` to allow Modifying/Updating existing tasks, rounding out the full CRUD operations.
