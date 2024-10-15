package main

import (
	"fmt"
	"net/http"

	"github.com/vaanskii/rssagg/internal/auth"
	"github.com/vaanskii/rssagg/internal/database"
)

type AuthHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *ApiConfig) MiddlewareAuth(handler AuthHandler) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			RespondWithErorr(w, 403, fmt.Sprintf("Auth errorr: %v", err))
		}
	
		user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			RespondWithErorr(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
		}

		handler(w, r, user)
	}
}