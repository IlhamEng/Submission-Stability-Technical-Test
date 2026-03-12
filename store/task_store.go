package store

import (
	"stability-test-task-api/models"
	"sync"
)

var (
	// mu Melindungi slice Tasks dan nextID dari Data Race (Akses bersamaan)
	mu    sync.RWMutex
	Tasks = []models.Task{
		{ID: 1, Title: "Learn Go", Done: false},
		{ID: 2, Title: "Build API", Done: false},
	}
	nextID = 3 // ID auto-increment selanjutnya
)

// GetAllTasks mereturn salinan seluruh data task 
func GetAllTasks() []models.Task {
	mu.RLock() // RLock = Kunci hanya-baca (bisa dibaca berbarengan)
	defer mu.RUnlock()

	// mereturn kopian isi arraynya.
	tasksCopy := make([]models.Task, len(Tasks))
	copy(tasksCopy, Tasks)
	
	return tasksCopy
}

func GetTaskByID(id int) *models.Task {
	mu.RLock()
	defer mu.RUnlock()

	for _, t := range Tasks {
		if t.ID == id {
			// Kita return nilai fix salinannya
			taskCopy := t
			return &taskCopy
		}
	}
	return nil
}

// AddTask menangani generate ID otomatis & Mutex penulisan data
// Membalikkan (return) pointer hasil task yang sudah terbuat ID-nya
func AddTask(task models.Task) models.Task {
	mu.Lock() 
	defer mu.Unlock()

	task.ID = nextID
	nextID++
	
	Tasks = append(Tasks, task)
	return task
}

// UpdateTask memperbarui data task berdasarkan ID
func UpdateTask(id int, updatedData models.Task) *models.Task {
	mu.Lock()
	defer mu.Unlock()

	for i, t := range Tasks {
		if t.ID == id {
			Tasks[i].Title = updatedData.Title
			Tasks[i].Done = updatedData.Done
			
			taskCopy := Tasks[i]
			return &taskCopy
		}
	}
	
	return nil
}

// DeleteTask mereturn boolean; true jika berhasil dihapus, false jika id tidak ketemu
func DeleteTask(id int) bool {
	mu.Lock()
	defer mu.Unlock()

	for i, t := range Tasks {
		if t.ID == id {
			Tasks = append(Tasks[:i], Tasks[i+1:]...)
			return true // Menghentikan iterasi setelah dihapus
		}
	}
	
	return false
}
