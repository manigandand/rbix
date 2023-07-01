package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func main() {
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
		middleware.Recoverer,
	)

	router.HandleFunc("/{container_id}/ws", angagoServer)

	log.Println("angago listening on ", Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", Port), router); err != nil {
		log.Fatal(err.Error())
	}
}
