package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/vaanskii/rssagg/internal/database"

	_ "github.com/lib/pq"
)

type ApiConfig struct {
	DB *database.Queries
}

func main() {
	// Load environment variables from .env file
	godotenv.Load()

	// Check if the port avaliable in environment
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment.")
	}

	// Check if the DB_URL avaliable in environment
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL is not found in the environment.")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Can't connect to database: ", err)
	}

	queries := database.New(conn)

	apiCfg := ApiConfig{
		DB: queries,
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
	v1Router.Post("/users", apiCfg.CreateUserHandler) // handles the creation of a new user in the system.
	v1Router.Get("/users", apiCfg.MiddlewareAuth(apiCfg.HandlerGetUserByAPIKey))
	v1Router.Post("/feeds", apiCfg.MiddlewareAuth(apiCfg.HandlerCreateFeed))

	// Mount the v1Router under the "/v1" path in the main router
	router.Mount("/v1", v1Router)

	// Create a HTTP Server
	srv := &http.Server{
		Handler: router,
		Addr: ":" + portString,
	}

	// Start the server and log the port
	log.Printf("Server starting on port: %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}