package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {

	maxBytes := 1_048_576 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

func (app *app) writeJSON(w http.ResponseWriter, status int, data any) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	return encoder.Encode(data)
}

func (app *app) writeJSONErr(w http.ResponseWriter, status int, data any) error {
	type errEnv struct {
		Error any `json:"error"`
	}
	return app.writeJSON(w, status, &errEnv{Error: data})
}

func (app *app) jsonResponse(w http.ResponseWriter, status int, data any) error {
	type jsonEnv struct {
		Data any `json:"data"`
	}
	return app.writeJSON(w, status, &jsonEnv{Data: data})
}
