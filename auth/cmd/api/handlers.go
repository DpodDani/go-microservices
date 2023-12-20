package main

import (
	"errors"
	"fmt"
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
		toolbox.ErrorJson(w, err, http.StatusBadRequest)
	}

	// validate user against database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		toolbox.ErrorJson(w, errors.New("invalid credentials"), http.StatusBadRequest)
	}

	matches, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !matches {
		toolbox.ErrorJson(w, errors.New("invalid credentials"), http.StatusUnauthorized)
	}

	response := toolbox.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s! 🎉", user.Email),
		Data:    user,
	}

	toolbox.WriteJson(w, http.StatusAccepted, response, nil)
}
