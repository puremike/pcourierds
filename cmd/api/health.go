package main

import "net/http"

type healthResponse struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Message     string `json:"message"`
	AppVersion  string `json:"app_version"`
}

// HealthCheck godoc
//
//	 @Summary		Get health
//		@Description	Returns the status of the application
//		@Tags			health
//		@Accept			json
//		@Produce		json
//		@Success		200	{object}	healthResponse
//		@Router			/health [get]
func (app *application) health(w http.ResponseWriter, r *http.Request) {

	healthStr := healthResponse{
		Status:      "Ok",
		Environment: app.config.env,
		Message:     "CDS Application is healthy",
		AppVersion:  "1.0.0",
	}

	err := app.jsonResponse(w, http.StatusOK, healthStr)
	if err != nil {
		app.internalServer(w, r, err)
		return
	}
}
