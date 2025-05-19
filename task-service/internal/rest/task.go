package rest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/dielit66/task-management-system/internal/entities"
	"github.com/dielit66/task-management-system/internal/logger"
	"github.com/dielit66/task-management-system/internal/middleware"
	"github.com/gorilla/mux"
)

type TaskUseCase interface {
	GetAllById(ctx context.Context, id int) ([]*entities.Task, error)
	GetById(ctx context.Context, id int) (*entities.Task, error)
	Create(ctx context.Context, t *entities.CreateTaskDto) error
	Update(ctx context.Context, t *entities.Task) error
	Delete(ctx context.Context, id int) error
}

type TaskHandler struct {
	Usecase TaskUseCase
	logger  logger.ILogger
}

func NewTaskHandler(m *mux.Router, uc TaskUseCase, l logger.ILogger) {
	handler := TaskHandler{
		Usecase: uc,
		logger:  l,
	}

	m.HandleFunc("/tasks", handler.GetAllByUserId).Methods("GET")
	m.HandleFunc("/tasks", handler.Create).Methods("POST")
	m.HandleFunc("/tasks/{id:[0-9]+}", handler.GetById).Methods("GET")
	m.HandleFunc("/tasks/{id:[0-9]+}", handler.Update).Methods("PUT")
	m.HandleFunc("/tasks/{id:[0-9]+}", handler.Delete).Methods("DELETE")

	m.Use(middleware.JwtPayloadMiddleware(l))
}

func (h *TaskHandler) GetById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid task ID", "invalid_id")
		return
	}

	task, err := h.Usecase.GetById(context.Background(), id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "Task not found", "not_found")
		return
	}

	body, err := json.Marshal(task)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Error marshalling response", "marshal_error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (h *TaskHandler) GetAllByUserId(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		h.writeError(w, http.StatusUnauthorized, "User not authenticated", "unauthorized")
		return
	}

	h.logger.Debug("Fetching tasks for user", "user_id", userID)
	tasks, err := h.Usecase.GetAllById(context.Background(), userID)
	if err != nil {
		h.logger.Error("Failed to fetch tasks", "user_id", userID, "error", err)
		h.writeError(w, http.StatusInternalServerError, "Failed to fetch tasks", "fetch_error")
		return
	}

	body, err := json.Marshal(tasks)
	if err != nil {
		h.logger.Error("Error marshalling response", "user_id", userID, "error", err)
		h.writeError(w, http.StatusInternalServerError, "Error marshalling response", "marshal_error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		h.logger.Warn("Invalid userID in context", "user_id", userID)
		h.writeError(w, http.StatusUnauthorized, "User not authenticated", "unauthorized")
		return
	}

	h.logger.Debug("Creating task for user", "user_id", userID)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Error reading request body", "read_error")
		return
	}

	var dto entities.CreateTaskDto
	err = json.Unmarshal(body, &dto)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Error parsing request body", "parse_error")
		return
	}

	if dto.Title == "" {
		h.writeError(w, http.StatusBadRequest, "Title is required", "missing_title")
		return
	}

	dto.UserID = userID

	err = h.Usecase.Create(context.Background(), &dto)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to create task", "create_error")
		return
	}

	h.logger.Info("Task created", "user_id", userID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("id"))
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid task ID", "invalid_id")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Error reading request body", "read_error")
		return
	}

	var task entities.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Error parsing request body", "parse_error")
		return
	}

	if task.Title == "" {
		h.writeError(w, http.StatusBadRequest, "Title is required", "missing_title")
		return
	}

	task.ID = id
	err = h.Usecase.Update(context.Background(), &task)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "Task not found", "not_found")
		return
	}

	h.logger.Info("Task updated", "task_id", id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid task ID", "invalid_id")
		return
	}

	err = h.Usecase.Delete(context.Background(), id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "Task not found", "not_found")
		return
	}

	h.logger.Info("Task deleted", "task_id", id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}

func (h *TaskHandler) writeError(w http.ResponseWriter, status int, message, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
		Code:  code,
	})
	h.logger.Error("Request failed", "status", status, "message", message, "code", code)
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}
