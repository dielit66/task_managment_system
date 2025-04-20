package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/dielit66/task-management-system/internal/entities"
	"github.com/dielit66/task-management-system/internal/usecases"
)

type UserHandler struct {
	usecase *usecases.UserUseCase
}

func NewUserHandler(uc *usecases.UserUseCase) *UserHandler {
	return &UserHandler{
		usecase: uc,
	}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error while reading body"))
		return
	}

	var user entities.User

	err = json.Unmarshal(body, &user)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error while unmarshalling json"))
		return
	}

	err = h.usecase.RegisterUser(context.Background(), user.Username, user.Email, user.PasswordHash)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error while creating user"))
		return
	}

	w.WriteHeader(201)
	w.Write([]byte{1})
}
