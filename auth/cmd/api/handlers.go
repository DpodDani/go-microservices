package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	toolbox "github.com/DpodDani/go-microservices-toolbox/json"
)

func logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	request, err := http.NewRequest(http.MethodPost, "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}

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

	// log authentication
	// note: for some reason, in the tutorial we return from this function if
	// logRequest() call returns an error... Decided not worth doing that
	err = logRequest("auth", fmt.Sprintf("%s logged in", requestPayload.Email))
	if err != nil {
		log.Printf("😔 - Failed to log auth for email: %s: %s", requestPayload.Email, err)
	}

	response := toolbox.JsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s! 🎉", user.Email),
		Data:    user,
	}

	toolbox.WriteJson(w, http.StatusAccepted, response, nil)
}
