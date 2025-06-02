package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/puremike/pcourierds/internal/models"
)

type dispatcherApplyRequest struct {
	VehicleType        string `json:"vehicle_type" binding:"required,oneof=car motorcycle"`
	VehiclePlateNumber string `json:"vehicle_plate_number" binding:"required"`
	VehicleYear        int    `json:"vehicle_year" binding:"required"`
	VehicleModel       string `json:"vehicle_model" binding:"required"`
	DriverLicense      string `json:"driver_license" binding:"required"`
}

type dispatcherResponse struct {
	ID                 string `json:"id"`
	UserID             string `json:"user_id"`
	VehicleType        string `json:"vehicle_type"`
	VehiclePlateNumber string `json:"vehicle_plate_number"`
	VehicleYear        int    `json:"vehicle_year"`
	VehicleModel       string `json:"vehicle_model"`
	DriverLicense      string `json:"driver_license"`
	Status             string `json:"status"` // pending, approved, rejected
	CreatedAt          string `json:"created_at"`
}

// CreateDispatcherApplication godoc
//
//	@Summary		Create dispatcher application
//	@Description	Apply as a dispatcher
//	@Tags			DispatchersApply
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dispatcherApplyRequest	true	"Dispatcher Application payload"
//	@Success		200		{object}	dispatcherResponse
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/dispatchers/apply [post]
//
//	@Security		BearerAuth
func (app *application) dispatcherApply(c *gin.Context) {
	var payload dispatcherApplyRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authUser := app.getUserFromContext(c)

	if authUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if authUser.Role == "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "only users can create applications"})
		return
	}

	apply := &models.DispatcherApplication{
		UserID:             authUser.ID,
		VehicleType:        payload.VehicleType,
		VehiclePlateNumber: payload.VehiclePlateNumber,
		VehicleYear:        payload.VehicleYear,
		VehicleModel:       payload.VehicleModel,
		DriverLicense:      payload.DriverLicense,
		Status:             "pending",
	}

	submittedApplylication, err := app.store.DispatcherApplications.DispatcherApplication(c.Request.Context(), apply)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create dispatcher application"})
		return
	}

	c.JSON(http.StatusCreated, dispatcherResponse{
		ID:                 submittedApplylication.ID,
		UserID:             submittedApplylication.UserID,
		VehicleType:        submittedApplylication.VehicleType,
		VehiclePlateNumber: submittedApplylication.VehiclePlateNumber,
		VehicleYear:        submittedApplylication.VehicleYear,
		VehicleModel:       submittedApplylication.VehicleModel,
		DriverLicense:      submittedApplylication.DriverLicense,
		Status:             submittedApplylication.Status,
		CreatedAt:          submittedApplylication.CreatedAt.Format(time.RFC3339),
	})
}

// GetDispatherApplications godoc
//
//	@Summary		Get Dispatcher Applications
//	@Description	Get All Dispatcher Applications
//	@Tags			DispatchersApply
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dispatcherResponse
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/admin/dispatcher-applications [get]
//
// @Security		BearerAuth
func (app *application) getAllApplications(c *gin.Context) {

	authUser := app.getUserFromContext(c)
	if authUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	applications, err := app.store.DispatcherApplications.GetAllApplications(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve applications"})
		return
	}

	var response []dispatcherResponse
	for _, application := range *applications {
		response = append(response, dispatcherResponse{
			ID:                 application.ID,
			UserID:             application.UserID,
			VehicleType:        application.VehicleType,
			VehiclePlateNumber: application.VehiclePlateNumber,
			VehicleYear:        application.VehicleYear,
			VehicleModel:       application.VehicleModel,
			DriverLicense:      application.DriverLicense,
			Status:             application.Status,
			CreatedAt:          application.CreatedAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, response)
}

// GetDispatcherById godoc
//
//	@Summary		Get Dispatcher Application
//	@Description	Get Dispatcher Application by ID
//	@Tags			DispatchersApply
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Dispatcher Application ID"
//	@Success		200	{object}	dispatcherResponse
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/admin/dispatcher-applications/{id} [get]
//
// @Security		BearerAuth
func (app *application) getDispatcherApplicationById(c *gin.Context) {

	authUser := app.getUserFromContext(c)
	if authUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	dispatcherApp := app.getDispatcherAppFromContext(c)
	if dispatcherApp == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "dispatcher application not found"})
		return
	}

	c.JSON(http.StatusOK, dispatcherResponse{
		ID:                 dispatcherApp.ID,
		UserID:             dispatcherApp.UserID,
		VehicleType:        dispatcherApp.VehicleType,
		VehiclePlateNumber: dispatcherApp.VehiclePlateNumber,
		VehicleYear:        dispatcherApp.VehicleYear,
		VehicleModel:       dispatcherApp.VehicleModel,
		DriverLicense:      dispatcherApp.DriverLicense,
		Status:             dispatcherApp.Status,
		CreatedAt:          dispatcherApp.CreatedAt.Format(time.RFC3339),
	})
}
