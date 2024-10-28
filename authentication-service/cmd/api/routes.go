package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// routes sets up the application's HTTP routes and middleware
func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	// Allow CORS from any origin with specified methods and headers
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Allow these HTTP methods
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}, // Allow these headers
		ExposedHeaders:   []string{"Link"}, // Expose these headers
		AllowCredentials: true, // Allow credentials
		MaxAge:           300,  // Cache preflight response for 5 minutes
	}))

	// Add a heartbeat endpoint for health checks
	mux.Use(middleware.Heartbeat("/ping"))

	// Route for user authentication
	mux.Post("/authenticate", app.Authenticate)

	return mux
}
