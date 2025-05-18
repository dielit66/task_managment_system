package rest

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/dielit66/task-management-system/internal/entities"
	app "github.com/dielit66/task-management-system/internal/errors"
	"github.com/dielit66/task-management-system/internal/logger"
	"github.com/gorilla/mux"
)

type UserService interface {
	RegisterUser(ctx context.Context, username string, email string, password string) error
	GetUser(ctx context.Context, id int) (*entities.User, error)
}

type UserHandler struct {
	Service UserService
	logger  logger.ILogger
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

func NewUserHandler(m *mux.Router, svc UserService, l logger.ILogger) {
	handler := &UserHandler{
		Service: svc,
		logger:  l,
	}

	m.HandleFunc("/users/register", handler.RegisterUser).Methods("POST")
	m.HandleFunc("/users/{id:[0-9]+}", handler.GetUser).Methods("GET")
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Error while reading body", "unknown")
		return
	}

	var user entities.CreateUserDto

	err = json.Unmarshal(body, &user)

	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Error while unmarshalling error", "unknown")
		return
	}

	err = h.Service.RegisterUser(context.Background(), user.Username, user.Email, user.Password)

	if err != nil {
		var appErr *app.AppError

		if errors.As(err, &appErr) {
			switch appErr.Type {
			case app.ErrInternal:
				h.logger.Error("Failed to create user", "username", user.Username, "error", err.Error())
				h.writeError(w, http.StatusInternalServerError, "Failed to create user", string(appErr.Type))
				return
			default:
				h.logger.Error("Unexpected error", "username", user.Username, "error", err.Error())
				h.writeError(w, http.StatusInternalServerError, "Unexpected error", string(appErr.Type))
				return
			}

		}
		h.logger.Error("Unexpected error", "username", user.Username, "error", err.Error())
		h.writeError(w, http.StatusInternalServerError, "internal server error", "unknown")
		return
	}

	h.logger.Info("User was created", "username", user.Username, "email", user.Email)
	w.WriteHeader(201)
	w.Write([]byte("1"))
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	if id == "" {
		w.WriteHeader(400)
		w.Write([]byte("Id is not provided"))
		return
	}

	idInt, err := strconv.Atoi(id)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("id is not a number"))
		return
	}

	user, err := h.Service.GetUser(context.Background(), idInt)

	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(400)
		w.Write([]byte("User is not found"))
		return
	}

	body, err := json.Marshal(user)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error while marshalling json"))
		return
	}

	w.WriteHeader(200)
	w.Write(body)
}

func (h *UserHandler) writeError(w http.ResponseWriter, status int, message, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
		Code:  code,
	})
}
