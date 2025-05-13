package main

import "net/http"

func (app *app) health(w http.ResponseWriter, r *http.Request) {

	healthStr := map[string]string{
		"Status":      "Ok",
		"Environment": app.config.env,
		"Message":     "CDS Application is healthy",
		"App_Version": "1.0.0",
	}

	err := app.jsonResponse(w, http.StatusOK, healthStr)
	if err != nil {
		app.internalServer(w, r, err)
		return
	}
}
