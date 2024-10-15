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
func (apiCfg *ApiConfig) HandlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		URL string 	`json:"url"`
	}

	decoder := json.NewDecoder(r.Body) // Create a JSON decoder for the request body

	params := parameters{} // Initialize parameters struct
	err := decoder.Decode(&params) // Decode the JSON into the parameters struct
	if err != nil {
		RespondWithErorr(w, 400, fmt.Sprintf("Error parsing JSON %s", err))
		return
	}

	// Create a new user in the database
	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:     params.Name,
		Url: params.URL,
		UserID: user.ID,
	})
	if err != nil {
		RespondWithErorr(w, 400, fmt.Sprintf("Couldn't create feed: %s", err)) // Respond with an error if user creation fails
		return
	}

	// Respond with the newly created user in JSON format
	RespondWithJSON(w, 201, DatabaseFeedToFeed(feed))
}
