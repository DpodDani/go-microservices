package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/DpodDani/go-microservices-toolbox/json"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.ReadJson(w, r, &requestPayload)
	if err != nil {
		json.ErrorJson(w, err, http.StatusBadRequest)
	}

	// validate user against database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		json.ErrorJson(w, errors.New("invalid credentials"), http.StatusBadRequest)
	}

	matches, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !matches {
		json.ErrorJson(w, errors.New("invalid credentials"), http.StatusUnauthorized)
	}

	response := json.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s! 🎉", user.Email),
		Data:    user,
	}

	err = json.WriteJson(w, http.StatusAccepted, response, nil)
}
