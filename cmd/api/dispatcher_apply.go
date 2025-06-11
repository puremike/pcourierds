package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/puremike/pcourierds/internal/models"
	"github.com/puremike/pcourierds/internal/store"
)

type dispatcherApplyRequest struct {
	VehicleType        string `json:"vehicle_type" binding:"required,oneof=car motorcycle"`
	VehiclePlateNumber string `json:"vehicle_plate_number" binding:"required"`
	VehicleYear        int    `json:"vehicle_year" binding:"required"`
	VehicleModel       string `json:"vehicle_model" binding:"required"`
	DriverLicense      string `json:"driver_license" binding:"required"`
}

type dispatcherAppResponse struct {
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

// type dispatcherResponse struct {
// 	ID                 string    `json:"id"`
// 	UserID             string    `json:"user_id"`
// 	ApplicationID      string    `json:"application_id"`
// 	VehicleType        string    `json:"vehicle_type"`
// 	VehiclePlateNumber string    `json:"vehicle_plate_number"`
// 	VehicleYear        int       `json:"vehicle_year"`
// 	VehicleModel       string    `json:"vehicle_model"`
// 	DriverLicense      string    `json:"driver_license"`
// 	ApprovedAt         time.Time `json:"approved_at"`
// 	IsActive           bool      `json:"isactive"`
// 	Rating             float32   `json:"rating"`
// 	CreatedAt          string    `json:"created_at"`
// }

// CreateDispatcherApplication godoc
//
//	@Summary		Create dispatcher application
//	@Description	Apply as a dispatcher
//	@Tags			DispatchersApply
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dispatcherApplyRequest	true	"Dispatcher Application payload"
//	@Success		201		{object}	dispatcherAppResponse
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

	authUser, err := app.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if authUser.Role == "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "only users can create applications"})
		return
	}

	existingApp, err := app.store.DispatcherApplications.GetApplicationByUserId(c.Request.Context(), authUser.ID)
	if err != nil && !errors.Is(err, store.ErrDispatcherApplicationNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check existing application"})
		return
	}
	if existingApp != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user already has an application"})
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

	c.JSON(http.StatusCreated, dispatcherAppResponse{
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
//	@Success		200	{object}	dispatcherAppResponse
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/admin/dispatcher-applications [get]
//
//	@Security		BearerAuth
func (app *application) getAllApplications(c *gin.Context) {

	authUser, err := app.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if authUser.Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	applications, err := app.store.DispatcherApplications.GetAllApplications(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve applications"})
		return
	}

	var response []dispatcherAppResponse
	for _, application := range *applications {
		response = append(response, dispatcherAppResponse{
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
//	@Success		200	{object}	dispatcherAppResponse
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/admin/dispatcher-applications/{id} [get]
//
//	@Security		BearerAuth
func (app *application) getDispatcherApplicationById(c *gin.Context) {

	authUser, err := app.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if authUser.Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	dispatcherApp, err := app.getDispatcherAppFromContext(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "dispatcher application not found"})
		return
	}

	c.JSON(http.StatusOK, dispatcherAppResponse{
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

// ApproveOrDenyDispatcherApplication godoc
//
//	@Summary		Approve or Deny a dispatcher application
//	@Description	Approve or Deny a dispatcher
//	@Tags			DispatchersApply
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		string				true	"Dispatcher Application ID"
//
//	@Success		200		{object}	map[string]string	"success: application rejected"
//	@Success		201		{object}	map[string]string	"success: application approved"
//
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/admin/approve-dispatcher/{userID} [patch]
//
//	@Security		BearerAuth
func (app *application) approveDenyApplication(c *gin.Context) {

	authUser, err := app.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if authUser.Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	dispatcherApp, err := app.getDispatcherAppByUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "dispatcher application not found"})
		return
	}

	// Check if application is already approved
	if dispatcherApp.Status == "approved" {
		c.JSON(http.StatusOK, gin.H{"message": "Dispatch Application Already Approved"})
		return
	}

	// Reject application based on the following logic

	if dispatcherApp.Status != "pending" || len(dispatcherApp.VehiclePlateNumber) != 8 || len(dispatcherApp.DriverLicense) != 12 || dispatcherApp.VehicleYear < 2008 {
		if err := app.store.DispatcherApplications.DeleteApplicationByUserId(c.Request.Context(), dispatcherApp.UserID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete dispatcher application"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Dispatch Application Rejected"})
		return
	}

	// Create and approve dispatch application
	dispatcher := &models.Dispatcher{
		UserID:             dispatcherApp.UserID,
		ApplicationID:      dispatcherApp.ID,
		VehicleType:        dispatcherApp.VehicleType,
		VehiclePlateNumber: dispatcherApp.VehiclePlateNumber,
		VehicleYear:        dispatcherApp.VehicleYear,
		VehicleModel:       dispatcherApp.VehicleModel,
		DriverLicense:      dispatcherApp.DriverLicense,
		ApprovedAt:         time.Now(),
		IsActive:           true,
		Rating:             0,
	}

	if err := app.store.Dispatchers.CreateDispatcher(c.Request.Context(), dispatcher); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create dispatcher"})
		return
	}

	dispatcherApp.Status = "approved"

	if err := app.store.DispatcherApplications.UpdateDispatchApplicationStatus(c.Request.Context(), dispatcherApp, dispatcherApp.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update dispatcher application"})
		return
	}

	user, err := app.store.Users.GetUserById(c.Request.Context(), dispatcherApp.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user"})
		return
	}

	user.Role = "dispatcher"
	if _, err := app.store.Users.UpdateUser(c.Request.Context(), user, user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	c.JSON(http.StatusCreated, "Your Dispatch Application has Been Approved")
}
