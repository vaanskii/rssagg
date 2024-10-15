package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Returning error message
func RespondWithErorr(w http.ResponseWriter, code int, msg string) {
	// checking if error code is client side error
	if code > 499 {
		log.Println("Responding with 500 erorr: ", msg)
	}
	// Declaring errorResponse struct
	type errResponse struct {
		Error string `json:"error"`
	}

	RespondWithJSON(w, code, errResponse{
		Error: msg,
	})
}
// Func which returning JSON format code
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	// Marshal the payload into JSON format.
	dat, err := json.Marshal(payload)
	if err != nil {
		// If marshaling fails, log the error and respond with a 500 Internal Server Error.
		log.Printf("Failed to marshal json response: %v", payload)
		w.WriteHeader(http.StatusInternalServerError) // status 500
		return
	}
	w.Header().Add("Content-Type", "application/json") // Set the Content-Type header to application/json.
	w.WriteHeader(code) // Write the HTTP status code to the response.
	w.Write(dat) // Write the JSON data to the response body.
}
