package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alokxcode/students-api/internal/config"
)

func main() {
	// load config
	cfg := config.MustLoad()
	// database setup
	// setup router

	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Students API!"))
	})

	// setup server

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	fmt.Println("server started")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Server did not Started", err)
	}
}
