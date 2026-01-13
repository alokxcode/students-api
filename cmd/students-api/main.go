package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alokxcode/students-api/internal/config"
	"github.com/alokxcode/students-api/internal/http/handlers/student"
	"github.com/alokxcode/students-api/internal/http/storage/sqlite"
)

func main() {
	// load config
	cfg := config.MustLoad()
	// database setup

	_, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("database initialised", slog.String("ENV", cfg.Env))

	// setup router

	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New())

	// setup server

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	slog.Info("Server started", slog.String("Address", cfg.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("Server did not Started", err)
		}

	}()

	<-done

	slog.Info("Shutting down the server")

	cxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(cxt); err != nil {
		slog.Error("failed to shutdown server, err:", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")

}
