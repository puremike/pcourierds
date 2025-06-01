package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/puremike/pcourierds/internal/models"
)

type dispatcherApplyRequest struct {
	Vehicle string `json:"vehicle" binding:"required"`
	License string `json:"license" binding:"required"`
	// Status  string `json:"status" binding:"required,oneof=pending approved rejected"`
}

type dispatcherResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Vehicle   string `json:"vehicle"`
	License   string `json:"license"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
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
		UserID:  authUser.ID,
		Vehicle: payload.Vehicle,
		License: payload.License,
		Status:  "pending",
	}

	applyApp, err := app.store.DispatcherApplications.DispatcherApplication(c.Request.Context(), apply)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create dispatcher application"})
		return
	}

	c.JSON(http.StatusCreated, dispatcherResponse{
		ID:        applyApp.ID,
		UserID:    applyApp.UserID,
		Vehicle:   applyApp.Vehicle,
		License:   applyApp.License,
		Status:    applyApp.Status,
		CreatedAt: applyApp.CreatedAt.Format(time.RFC3339),
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
			ID:        application.ID,
			UserID:    application.UserID,
			Vehicle:   application.Vehicle,
			License:   application.License,
			Status:    application.Status,
			CreatedAt: application.CreatedAt.Format(time.RFC3339),
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
		ID:        dispatcherApp.ID,
		UserID:    dispatcherApp.UserID,
		Vehicle:   dispatcherApp.Vehicle,
		License:   dispatcherApp.License,
		Status:    dispatcherApp.Status,
		CreatedAt: dispatcherApp.CreatedAt.Format(time.RFC3339),
	})
}
