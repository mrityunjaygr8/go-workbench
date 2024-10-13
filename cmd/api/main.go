package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

var healthy int32

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func healthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	if atomic.LoadInt32(&healthy) == 1 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
}

func longRequestHandler(w http.ResponseWriter, _ *http.Request) {
	time.Sleep(20 * time.Second)
	fmt.Fprintf(w, "Late current")
}

func main() {
	logger := log.New(os.Stdout, "http: ", log.LstdFlags)
	logger.Println("Server is starting...")
	logger.Println("Server Running on 8080")

	store, err := NewInMemoryStore()
	if err != nil {
		logger.Fatal(err)
	}

	app := &application{
		logger: logger,
		store:  store,
	}
	mux := app.mountRoutes()
	logger.Fatal(app.run(mux))

}
