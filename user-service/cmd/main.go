package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/dielit66/task-management-system/internal/handlers"
	repository "github.com/dielit66/task-management-system/internal/repository/postgres"
	"github.com/dielit66/task-management-system/internal/usecases"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://user:password@localhost:5432/task_management?sslmode=disable")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	defer db.Close()

	repo := repository.NewPostgresUserRepostiry(db)

	usecase := usecases.NewUserUseCase(repo)

	handler := handlers.NewUserHandler(usecase)

	router := mux.NewRouter()
	router.HandleFunc("/users/register", handler.RegisterUser).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))

}
