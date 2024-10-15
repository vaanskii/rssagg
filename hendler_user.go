package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/vaanskii/rssagg/internal/database"
)

// CreateUserHandler handles the creation of a new user in the system.
func (apiCfg *ApiConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body) // Create a JSON decoder for the request body

	params := parameters{} // Initialize parameters struct
	err := decoder.Decode(&params) // Decode the JSON into the parameters struct
	if err != nil {
		RespondWithErorr(w, 400, fmt.Sprintf("Error parsing JSON %s", err))
		return
	}

	// Create a new user in the database
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:     params.Name,
	})
	if err != nil {
		RespondWithErorr(w, 400, fmt.Sprintf("Couldn't create user: %s", err)) // Respond with an error if user creation fails
		return
	}

	// Respond with the newly created user in JSON format
	RespondWithJSON(w, 201, DatabaseUserToUser(user))
}


func (apiCfg *ApiConfig) HandlerGetUserByAPIKey(w http.ResponseWriter, r *http.Request, user database.User) {
	RespondWithJSON(w, 200, DatabaseUserToUser(user))
}