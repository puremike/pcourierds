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
