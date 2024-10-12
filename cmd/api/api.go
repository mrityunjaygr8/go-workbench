package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"
)

type application struct {
	logger *log.Logger
	store  store
}

type config struct {
}

func (a *application) mountRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", helloHandler)
	mux.HandleFunc("GET /healthz", healthCheckHandler)
	mux.HandleFunc("GET /long", longRequestHandler)
	mux.HandleFunc("POST /api/v1/users", a.createUserHandler)
	mux.HandleFunc("GET /api/v1/users", a.listUsersHandler)
	return mux
}

func (a *application) run(mux http.Handler) error {
	nextRequestID := func() string {
		return strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	server := &http.Server{
		Addr:         ":8080",
		Handler:      tracing(nextRequestID)(logging(a.logger)(mux)),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	atomic.StoreInt32(&healthy, 1)
	go func() {
		<-quit
		a.logger.Println("Server is shutting down...")
		atomic.StoreInt32(&healthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			a.logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
		}
		close(done)
	}()

	a.logger.Println("Server is ready to handle requests at :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		a.logger.Fatalf("Could not listen on :8080: %v\n", err)
	}
	<-done
	a.logger.Println("Server stopped")
	return nil
}
