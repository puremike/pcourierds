package main

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/puremike/pcourierds/internal/models"
)

func (app *application) getUserFromContext(c *gin.Context) (*models.User, error) {

	userContext, exists := c.Get("user")
	if !exists {
		return nil, errors.New("user not found in context")
	}

	user, ok := userContext.(*models.User)
	if !ok {
		return nil, errors.New("user context is not of type *models.User")
	}

	return user, nil
}

func (app *application) getDispatcherAppFromContext(c *gin.Context) (*models.DispatcherApplication, error) {

	dispatcherAppContext, exists := c.Get("dispatcherApp")
	if !exists {
		return nil, errors.New("dispatcherApp not found in context")
	}

	dispatcherApp, ok := dispatcherAppContext.(*models.DispatcherApplication)
	if !ok {
		return nil, errors.New("dispatcherApp context is not of type *models.DispatcherApplication")
	}

	return dispatcherApp, nil
}

func (app *application) getDispatcherAppByUserIdFromContext(c *gin.Context) (*models.DispatcherApplication, error) {

	dispatcherAppByUserIdContext, exists := c.Get("dispatcherAppByUserId")
	if !exists {
		return nil, errors.New("dispatcherAppByUserId not found in context")
	}

	dispatcherApp, ok := dispatcherAppByUserIdContext.(*models.DispatcherApplication)
	if !ok {
		return nil, errors.New("dispatcherAppByUserId context is not of type *models.DispatcherApplication")
	}

	return dispatcherApp, nil
}
