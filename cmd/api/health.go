package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type healthResponse struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Message     string `json:"message"`
	ApiVersion  string `json:"api_version"`
}

// HealthCheck godoc
//
//	@Summary		Get health
//	@Description	Returns the status of the application
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	healthResponse
//	@Router			/health [get]
func (app *application) health(c *gin.Context) {

	healthStr := healthResponse{
		Status:      "Ok",
		Environment: app.config.env,
		Message:     "Courier Delivery System is healthy",
		ApiVersion:  apiVersion,
	}

	c.JSON(http.StatusOK, healthStr)
}
