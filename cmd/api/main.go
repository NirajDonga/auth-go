package main

import (
	"go-auth/internal/httpserver"
	"log"
	"net/http"
	"time"
)

func main() {

	router := httpserver.NewRouter()

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Api running on: %s", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Printf("Server closed")
			return
		}
		log.Fatalf("Server error: %v", err)
	}

}
