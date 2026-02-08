package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"assignment2/internal/models"
)

type TaskHandler struct {
	mu     sync.Mutex
	tasks  map[int]models.Task
	nextID int
}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{
		tasks:  make(map[int]models.Task),
		nextID: 1,
	}
}

// GET /tasks
// GET /tasks?id=1
// GET /tasks?done=true   ✅ filtering (easy task)
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	h.mu.Lock()
	defer h.mu.Unlock()

	// get by id
	if idStr := r.URL.Query().Get("id"); idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid id"})
			return
		}

		task, ok := h.tasks[id]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "task not found"})
			return
		}

		json.NewEncoder(w).Encode(task)
		return
	}

	// filtering by done
	doneFilter := r.URL.Query().Get("done")
	result := []models.Task{}

	for _, task := range h.tasks {
		if doneFilter == "" {
			result = append(result, task)
			continue
		}

		done, err := strconv.ParseBool(doneFilter)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid done value"})
			return
		}

		if task.Done == done {
			result = append(result, task)
		}
	}

	json.NewEncoder(w).Encode(result)
}

// POST /tasks
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var body struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid title"})
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	task := models.Task{
		ID:    h.nextID,
		Title: body.Title,
		Done:  false,
	}

	h.tasks[h.nextID] = task
	h.nextID++

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// PATCH /tasks?id=1
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid id"})
		return
	}

	var body struct {
		Done *bool `json:"done"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Done == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid done value"})
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	task, ok := h.tasks[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "task not found"})
		return
	}

	task.Done = *body.Done
	h.tasks[id] = task

	json.NewEncoder(w).Encode(map[string]bool{"updated": true})
}
