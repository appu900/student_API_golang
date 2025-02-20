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

	"github.com/appu900/students-api/internal/config"
	student "github.com/appu900/students-api/internal/http/handler"
)

func main() {
	// load config
	cfg := config.MustLoad()

	// database setup
	// setup route

	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New())
	

	// setup server
	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("Server is running on", slog.String("address", cfg.Address))
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start server")
		}
	}()

	<-done

	slog.Info("shutting down the server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}
	slog.Info("server shutdown successfully")

}
