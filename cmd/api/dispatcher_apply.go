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
//	@Tags			Dispatchers
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dispatcherApplyRequest	true	"Dispatcher payload"
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
