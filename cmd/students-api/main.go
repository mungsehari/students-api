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

	"github.com/mungse.hari/students-api/internal/config"
)

func main() {
	// load configuration
	cfg := config.MustLoad()
	// database setup
	//setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Students API"))
	})
	//setup server

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}
	slog.Info("starting server on %s", slog.String("address", cfg.Addr))

	// start server
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("error starting server: %s", err.Error())
		}
	}()
	<-done

	// shutdown server
	slog.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("error shutting down server: %s", slog.String("error", err.Error()))
	}
	slog.Info("server shutdown successfully")

}
