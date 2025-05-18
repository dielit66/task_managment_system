package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dielit66/task-management-system/internal/auth"
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

	log.Println("Logger Initialization")

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

	l.Info("Creating jwt service")
	jwtService := auth.NewJWTService([]byte(cfg.Server.JWTSecret))

	l.Info("Creating new user usecase")
	usecase := usecases.NewAuthUseCase(repo, jwtService, l)

	l.Info("Creating router")
	router := mux.NewRouter()

	l.Info("Creating new user handler")
	rest.NewAuthHandler(router, usecase, l)

	port := fmt.Sprintf(":%s", cfg.Server.Port)

	srv := &http.Server{
		Addr:         port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		l.Info("Starting server", "port", port)
		if err := srv.ListenAndServe(); err != nil {
			l.Fatal(err.Error())
		}
	}()
	sdChan := make(chan os.Signal, 1)

	signal.Notify(sdChan, syscall.SIGINT, syscall.SIGTERM)

	<-sdChan
	l.Info("Shutting down the server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		l.Fatal("Server shutdown error", "error", err.Error())
	}

	l.Info("Server gracefully stopped")

}
