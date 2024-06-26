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

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/manigandand/adk/api"
	appmiddleware "github.com/manigandand/adk/middleware"
	"github.com/rs/cors"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Wrong length of arguments")
	}

	// initialize config
	Initialize(os.Args[1:]...)

	api.InitService("rbix-api", "v1.0.0")

	InitStore()
	InitOrchechestrator()

	router := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "OPTIONS", "DELETE"},
		AllowedHeaders: []string{
			"Origin", "Authorization", "Access-Control-Allow-Origin",
			"Access-Control-Allow-Header", "Accept",
			"Content-Type", "X-CSRF-Token",
		},
		ExposedHeaders: []string{
			"Content-Length", "Access-Control-Allow-Origin", "Origin",
		},
		AllowCredentials: true,
		MaxAge:           300,
	})

	// cross & loger middleware
	router.Use(cors.Handler)
	router.Use(
		middleware.Logger,
		appmiddleware.Recoverer,
	)
	router.Get("/", api.IndexHandeler)
	router.Get("/health", api.HealthHandeler)

	// routes
	router.Route("/v1", func(r chi.Router) {
		r.Method(http.MethodPost, "/try", api.Handler(newSqureXSessionHandler))
		r.Method(http.MethodGet, "/status/{container_id}", api.Handler(getContainerStatHandler))
		r.Method(http.MethodPost, "/stop/{termination_token}", api.Handler(stopContainerHandler))
	})

	server := http.Server{
		Addr:         fmt.Sprintf(":%s", Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	interruptChan := make(chan os.Signal, 1)
	go func() {
		signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		// Block until we receive our signal.
		<-interruptChan
		log.Println("os interrupt signal received, shutting down")

		// shutdown server
		if err := server.Shutdown(context.Background()); err != nil {
			log.Println("shtdown error:", err.Error())
		}

		os.Exit(1)
	}()

	log.Println("RbiX-api-server listening on ", Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", Port), router); err != nil {
		log.Fatal(err.Error())
	}
}
