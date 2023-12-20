package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	toolbox "github.com/DpodDani/go-microservices-toolbox/json"
)

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auto,omitempty"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := toolbox.JsonResponse{
		Error:   false,
		Message: "Hit the broker!",
	}

	_ = toolbox.WriteJson(w, http.StatusOK, payload, nil)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := toolbox.ReadJson(w, r, &requestPayload)
	if err != nil {
		toolbox.ErrorJson(w, err, http.StatusBadRequest)
	}

	switch requestPayload.Action {
	case "authenticate":
		app.authenticate(w, requestPayload.Auth)
	default:
		toolbox.ErrorJson(
			w,
			fmt.Errorf("unrecognised action: %s", requestPayload.Action),
			http.StatusBadRequest,
		)
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.MarshalIndent(a, "", "\t")
	request, err := http.NewRequest(http.MethodPost, "http://auth-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		toolbox.ErrorJson(w, err, http.StatusInternalServerError)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		toolbox.ErrorJson(w, err, http.StatusInternalServerError)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		toolbox.ErrorJson(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	if response.StatusCode != http.StatusOK {
		toolbox.ErrorJson(w, errors.New("error calling auth service"), http.StatusBadRequest)
		return
	}

	var authResponse toolbox.JsonResponse
	err = json.NewDecoder(response.Body).Decode(&authResponse)
	if err != nil {
		toolbox.ErrorJson(w, err, http.StatusBadRequest)
		return
	}

	if authResponse.Error {
		toolbox.ErrorJson(w, err, http.StatusUnauthorized)
		return
	}

	payload := toolbox.JsonResponse{
		Error:   false,
		Message: "Authenticated! 🎉",
		Data:    authResponse.Data,
	}
	toolbox.WriteJson(w, http.StatusAccepted, payload, nil)

}
