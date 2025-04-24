package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/dielit66/task-management-system/internal/entities"
	app "github.com/dielit66/task-management-system/internal/errors"
	"github.com/dielit66/task-management-system/internal/logger"
	"github.com/dielit66/task-management-system/internal/usecases"
)

type AuthHandler struct {
	usecase *usecases.AuthUseCase
	logger  logger.ILogger
}

func NewAuthHandler(uc *usecases.AuthUseCase, logger logger.ILogger) *AuthHandler {
	return &AuthHandler{
		usecase: uc,
		logger:  logger,
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

func (h *AuthHandler) writeError(w http.ResponseWriter, status int, message, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
		Code:  code,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	user := entities.AuthUserDto{}

	body, err := io.ReadAll(r.Body)

	defer r.Body.Close()

	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Error while reading body", "unknown")
		return
	}

	json.Unmarshal(body, &user)

	_, err = h.usecase.LoginUser(context.Background(), user.Username, user.Password)

	if err != nil {
		var appErr *app.AppError
		if errors.As(err, &appErr) {
			switch appErr.Type {
			case app.ErrNotFound:
				h.logger.Warn("User not found in handler", "username", user.Username)
				h.writeError(w, http.StatusNotFound, appErr.Message, string(appErr.Type))
				return
			case app.ErrUnauthorized:
				h.logger.Warn("Password missmatch", "username", user.Username)
				h.writeError(w, http.StatusUnauthorized, appErr.Message, string(appErr.Type))
				return
			default:
				h.logger.Error("Internal server error", "username", user.Username, "error", err.Error())
				h.writeError(w, http.StatusInternalServerError, "internal server error", string(appErr.Type))
				return
			}
		}
		h.logger.Error("Unexpected error", "username", user.Username, "error", err.Error())
		h.writeError(w, http.StatusInternalServerError, "internal server error", "unknown")
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("Login"))
}
