package main

import (
	"log"
	"net/http"

	toolbox "github.com/DpodDani/go-microservices-toolbox/json"
	"github.com/DpodDani/logger/data"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload JSONPayload

	_ = toolbox.ReadJson(w, r, &requestPayload)

	log.Printf("Writing log... %+v\n", requestPayload)

	entry := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}

	err := app.Models.LogEntry.Insert(entry)
	if err != nil {
		log.Println("Failed to insert log entry")
		toolbox.ErrorJson(w, err, http.StatusInternalServerError)
		return
	}

	resp := toolbox.JsonResponse{
		Error:   false,
		Message: "Logged message 🗣️",
	}

	toolbox.WriteJson(w, http.StatusAccepted, resp, nil)
}
