package json

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func ReadJson(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576 // 1MB limitation for JSON body

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must contain a single JSON object")
	}

	return nil
}

func WriteJson(w http.ResponseWriter, status int, data any, headers *http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if headers != nil {
		for key, value := range *headers {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func ErrorJson(w http.ResponseWriter, err error, statusCode int) error {
	var payload JsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return WriteJson(w, statusCode, payload, nil)
}
