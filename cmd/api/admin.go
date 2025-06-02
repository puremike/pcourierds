package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/puremike/pcourierds/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// CreateUserManually godoc
//
//	@Summary		Create user manually
//	@Description	Create a new user
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		createUserRequest	true	"User payload"
//	@Success		201		{object}	userResponse
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/admin/user [post]
//	@Security		BearerAuth
func (app *application) adminCreateUser(c *gin.Context) {

	var payload createUserRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if payload.Password != payload.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "passwords do not match"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	user := &models.User{
		Username: payload.Username,
		Email:    payload.Email,
		Role:     "user",
		Password: string(hashedPassword),
	}

	createdUser, err := app.store.Users.CreateUser(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, userResponse{
		ID:        createdUser.ID,
		Username:  createdUser.Username,
		Email:     createdUser.Email,
		Role:      createdUser.Role,
		CreatedAt: createdUser.CreatedAt.Format(time.RFC3339),
	})
}

// UpdateUserProfile godoc
//
//	@Summary		Update User Profile
//	@Description	Update Current User Profile
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Param			payload	body		userProfileUpdateRequest	true	"update credentials"
//	@Success		201		{object}	userResponse
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/admin/user/{id} [patch]
//
//	@Security		BearerAuth
func (app *application) adminUpdateProfile(c *gin.Context) {

	var payload userProfileUpdateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authUser := app.getUserFromContext(c)

	if authUser.Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userId := c.Param("id")

	existingUser, err := app.store.Users.GetUserById(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user"})
		return
	}

	user := &models.User{
		Role: "user",
	}

	if payload.Username != nil && *payload.Username != "" {
		user.Username = *payload.Username
	} else {
		user.Username = existingUser.Username
	}

	if payload.Email != nil && *payload.Email != "" {
		user.Email = *payload.Email
	} else {
		user.Email = existingUser.Email
	}

	updatedUser, err := app.store.Users.UpdateUser(c.Request.Context(), user, existingUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	c.JSON(http.StatusCreated, userResponse{
		ID:        updatedUser.ID,
		Username:  updatedUser.Username,
		Email:     updatedUser.Email,
		Role:      updatedUser.Role,
		CreatedAt: updatedUser.CreatedAt.Format(time.RFC3339),
	})
}

// GetuserById godoc
//
//	@Summary		Get User
//	@Description	Get User by ID
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	userResponse
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/admin/users/{id} [get]
//
// @Security		BearerAuth
func (app *application) getUserById(c *gin.Context) {

	authUser := app.getUserFromContext(c)

	if authUser.Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userId := c.Param("id")
	user, err := app.store.Users.GetUserById(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user"})
		return
	}

	c.JSON(http.StatusOK, userResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	})
}

// Getusersgodoc
//
//	@Summary		Get Users
//	@Description	Get All Users
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	userResponse
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/admin/users/ [get]
//
// @Security		BearerAuth
func (app *application) getUsers(c *gin.Context) {

	authUser := app.getUserFromContext(c)

	if authUser.Role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	users, err := app.store.Users.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve users"})
		return
	}

	var response []userResponse
	for _, user := range *users {
		response = append(response, userResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, response)
}
