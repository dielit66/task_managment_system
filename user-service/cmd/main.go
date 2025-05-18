package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dielit66/task-management-system/internal/config"
	"github.com/dielit66/task-management-system/internal/logger"
	repository "github.com/dielit66/task-management-system/internal/repository/postgres"
	"github.com/dielit66/task-management-system/internal/rest"
	"github.com/dielit66/task-management-system/internal/usecases"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	log.Println("Loading config")
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalln(err)
	}

	l, err := logger.NewZapLogger(cfg.LogLevel)

	if err != nil {
		log.Fatalf("Failed to initialize logger with level %d, error: %v", cfg.LogLevel, err)
	}

	dbDsn := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s host=%s", cfg.Databse.Username, cfg.Databse.Name, cfg.Databse.Password, cfg.Databse.Host)

	l.Debug("Trying to connect to database", "username", cfg.Databse.Username, "name", cfg.Databse.Name, "password", cfg.Databse.Password, "Host", cfg.Databse.Host)

	db, err := sqlx.Connect("postgres", dbDsn)

	if err != nil {
		l.Fatal("Error while connecting to databse", "err", err.Error())
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		l.Fatal("Db not respond on ping", "err", err.Error())
	} else {
		l.Info("Successfully connected to db")
	}

	l.Info("Creating new user repository")
	repo := repository.NewPostgresUserRepostiry(db, l)

	l.Info("Creating new user usecase")
	usecase := usecases.NewUserUseCase(repo, l)

	l.Info("Creating router")
	router := mux.NewRouter()

	l.Info("Creating new user handler")
	rest.NewUserHandler(router, usecase, l)

	port := fmt.Sprintf(":%s", cfg.Server.Port)

	l.Info("Starting server", "port", port)
	log.Fatal(http.ListenAndServe(port, router))

}
