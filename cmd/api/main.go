package main

import (
	"context"
	"go-auth/internal/app"
	"go-auth/internal/httpserver"
	"log"
	"net/http"
	"time"
)

func main() {

	ctx := context.Background()

	a, err := app.New(ctx)
	if err != nil {
		log.Fatalf("Startup Failed: %v", err)
	}

	defer func() {
		if err := a.Close(ctx); err != nil {
			log.Printf("Shut down warning: %v", err)
		}
	}()

	router := httpserver.NewRouter(a)

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
