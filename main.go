package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	godotenv.Load()

	// Check if the port avaliable in environment
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment.")
	}

	// Creating new chi router
	router := chi.NewRouter()

	// Set up CORS Middleware
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: false,
		MaxAge: 300,
	}))

	// Create a sub router for API Version 1
	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", HandlerReadiness) // Set up readiness handler for the '/ready' endpoint where allowed method is only GET
	v1Router.Get("/err", HandlerErr) // Handler endpoint for errors

	// Mount the v1Router under the "/v1" path in the main router
	router.Mount("/v1", v1Router)

	// Create a HTTP Server
	srv := &http.Server{
		Handler: router,
		Addr: ":" + portString,
	}

	// Start the server and log the port
	log.Printf("Server starting on port: %v", portString)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}