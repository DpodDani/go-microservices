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

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
	Mail   MailPayload `json:"mail,omitempty"`
}

func send_request(method string, url string, data []byte) (*http.Response, error) {
	request, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	if len(data) > 0 {
		request.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
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
	case "log":
		app.log(w, requestPayload.Log)
	case "mail":
		app.sendMail(w, requestPayload.Mail)
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
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		toolbox.ErrorJson(w, err, http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		toolbox.ErrorJson(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	if response.StatusCode != http.StatusAccepted {
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

func (app *Config) log(w http.ResponseWriter, entry LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	response, err := send_request(http.MethodPost, "http://logger-service/log", jsonData)
	if err != nil {
		log.Println("Failed to send request to logger service")
		toolbox.ErrorJson(w, err, http.StatusInternalServerError)
		return
	}

	if response.StatusCode != http.StatusAccepted {
		msg := "log service failed to process log entry"
		log.Println(msg)
		toolbox.ErrorJson(w, errors.New(msg), http.StatusBadRequest)
		return
	}

	var jsonResponse toolbox.JsonResponse
	jsonResponse.Error = false
	jsonResponse.Message = "Successfully processed log entry!"

	toolbox.WriteJson(w, http.StatusAccepted, jsonResponse, nil)
}

func (app *Config) sendMail(w http.ResponseWriter, mail MailPayload) {
	jsonData, _ := json.MarshalIndent(mail, "", "\t")
	response, err := send_request(http.MethodPost, "http://mail-service/send", jsonData)
	if err != nil {
		log.Println("Failed to send request to mail service")
		toolbox.ErrorJson(w, err, http.StatusInternalServerError)
		return
	}

	if response.StatusCode != http.StatusAccepted {
		msg := "mail service failed to process mail entry"
		log.Println(msg)
		toolbox.ErrorJson(w, errors.New(msg), http.StatusBadRequest)
		return
	}

	var jsonResponse toolbox.JsonResponse
	jsonResponse.Error = false
	jsonResponse.Message = fmt.Sprintf("Successfully sent mail to: %s!", mail.To)

	toolbox.WriteJson(w, http.StatusAccepted, jsonResponse, nil)
}
