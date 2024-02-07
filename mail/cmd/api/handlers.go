package main

import (
	"net/http"

	toolbox "github.com/DpodDani/go-microservices-toolbox/json"
)

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:"message"`
	}

	var requestPayload mailMessage

	err := toolbox.ReadJson(w, r, &requestPayload)
	if err != nil {
		toolbox.ErrorJson(w, err, http.StatusBadRequest)
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		toolbox.ErrorJson(w, err, http.StatusInternalServerError)
	}

	payload := toolbox.JsonResponse{
		Error:   false,
		Message: "sent to" + msg.To,
		Data:    msg.Data,
	}

	toolbox.WriteJson(w, http.StatusAccepted, payload, nil)

}
