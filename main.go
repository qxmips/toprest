package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/qxmips/toprest/resources"
)

func main() {

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "9090"
	}

	var bindAddress = "0.0.0.0:" + httpPort
	sm := mux.NewRouter()
	getRes := sm.Methods(http.MethodGet).Subrouter()
	getRes.HandleFunc("/resources", resources.GetResources)

	stop := make(chan struct{}, 1)

	s := &http.Server{
		Addr:         bindAddress,
		Handler:      sm,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// start the server
	go func() {
		log.Printf("Server listening at %s", bindAddress)
		err := s.ListenAndServe()
		// ListenAndServe always returns a non-nil error. After Shutdown or
		// Close, the returned error is ErrServerClosed
		if err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe ERROR: %v", err)
			close(stop)
			os.Exit(1)
		}

	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, os.Kill)
	sig := <-c
	log.Println("Recieved SIG:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s.Shutdown(ctx)

}
