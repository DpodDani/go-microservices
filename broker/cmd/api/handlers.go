package main

import (
	"net/http"

	toolbox "github.com/DpodDani/go-microservices-toolbox/json"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := toolbox.JsonResponse{
		Error:   false,
		Message: "Hit the broker!",
	}

	_ = toolbox.WriteJson(w, http.StatusOK, payload, nil)
}
