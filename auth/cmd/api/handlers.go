package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	toolbox "github.com/DpodDani/go-microservices-toolbox/json"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := toolbox.ReadJson(w, r, &requestPayload)
	if err != nil {
		log.Printf("❌ - Couldn't parse request")
		toolbox.ErrorJson(w, err, http.StatusBadRequest)
		return
	}

	// validate user against database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil || user == nil {
		log.Printf("❌ - Could not find email: %s\n", requestPayload.Email)
		toolbox.ErrorJson(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	matches, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !matches {
		log.Printf("❌ - Invalid password for email: %s\n", requestPayload.Email)
		toolbox.ErrorJson(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	response := toolbox.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s! 🎉", user.Email),
		Data:    user,
	}

	toolbox.WriteJson(w, http.StatusAccepted, response, nil)
}
