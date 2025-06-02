package main

import (
	"github.com/gin-gonic/gin"
	"github.com/puremike/pcourierds/internal/models"
)

func (app *application) getUserFromContext(c *gin.Context) *models.User {

	userContext, exists := c.Get("user")
	if !exists {
		return &models.User{}
	}

	user, ok := userContext.(*models.User)
	if !ok {
		return &models.User{}
	}

	return user
}

func (app *application) getDispatcherAppFromContext(c *gin.Context) *models.DispatcherApplication {

	dispatcherAppContext, exists := c.Get("dispatcherApp")
	if !exists {
		return &models.DispatcherApplication{}
	}

	dispatcherApp, ok := dispatcherAppContext.(*models.DispatcherApplication)
	if !ok {
		return &models.DispatcherApplication{}
	}

	return dispatcherApp
}

func (app *application) getDispatcherAppByUserIdFromContext(c *gin.Context) *models.DispatcherApplication {

	dispatcherAppByUserIdContext, exists := c.Get("dispatcherAppByUserId")
	if !exists {
		return &models.DispatcherApplication{}
	}

	dispatcherApp, ok := dispatcherAppByUserIdContext.(*models.DispatcherApplication)
	if !ok {
		return &models.DispatcherApplication{}
	}

	return dispatcherApp
}
