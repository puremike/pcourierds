package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/puremike/pcourierds/internal/models"
	"github.com/puremike/pcourierds/internal/store"
	"golang.org/x/crypto/bcrypt"
)

type createUserRequest struct {
	Username        string `json:"username" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,passwd"`
	ConfirmPassword string `json:"confirm_password" binding:"required,passwd"`
}

type userResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

type userProfileUpdateRequest struct {
	Username *string `json:"username"`
	Email    *string `json:"email" binding:"omitempty,email"`
}

type updatePasswordRequest struct {
	OldPassword     string `json:"old_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,passwd"`
	ConfirmPassword string `json:"confirm_password" binding:"required,passwd"`
}

// CreateUser godoc
//
//	@Summary		Create user
//	@Description	Create a new user
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		createUserRequest	true	"User payload"
//	@Success		201		{object}	userResponse
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/auth/signup [post]
func (app *application) createUser(c *gin.Context) {

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

// loginUser handles user login and returns a JWT token if credentials are valid.
//
//	@Summary		Login User
//	@Description	Authenticates a user using email and password, and returns a JWT token on success.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		loginRequest	true	"Login credentials"
//	@Success		200		{object}	loginResponse
//	@Failure		400		{object}	gin.H	"Bad Request - invalid input"
//	@Failure		401		{object}	gin.H	"Unauthorized - invalid credentials"
//	@Failure		500		{object}	gin.H	"Internal Server Error"
//	@Router			/auth/login [post]
func (app *application) login(c *gin.Context) {
	var payload loginRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := app.store.Users.GetUserByEmail(c.Request.Context(), payload.Email)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"iss":  app.config.authConfig.iss,
		"aud":  app.config.authConfig.aud,
		"iat":  time.Now().Unix(),
		"nbf":  time.Now().Unix(),
		"exp":  time.Now().Add(app.config.authConfig.tokenExp).Unix(),
	}

	token, err := app.jwtAuth.GenerateToken(claims)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.SetCookie("jwt", token, int(app.config.authConfig.tokenExp.Seconds()), "/", "localhost", false, true)
	c.SetSameSite(http.SameSiteLaxMode)

	res := loginResponse{ID: user.ID, Username: user.Username}
	c.JSON(http.StatusOK, res)
}

// GetLoggedUserProfile godoc
//
//	@Summary		Get User Profile
//	@Description	Get Current User Profile
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	userResponse
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/auth/me [get]
//
//	@Security		BearerAuth
func (app *application) userProfile(c *gin.Context) {
	user, err := app.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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

// UpdateUserProfile godoc
//
//	@Summary		Update User Profile
//	@Description	Update Current User Profile
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		userProfileUpdateRequest	true	"update credentials"
//	@Success		201		{object}	userResponse
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/auth/update-profile [patch]
//
//	@Security		BearerAuth
func (app *application) updateProfile(c *gin.Context) {

	var payload userProfileUpdateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authUser, err := app.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user := &models.User{
		Role: "user",
	}

	if payload.Username != nil && *payload.Username != "" {
		user.Username = *payload.Username
	} else {
		user.Username = authUser.Username
	}

	if payload.Email != nil && *payload.Email != "" {
		user.Email = *payload.Email
	} else {
		user.Email = authUser.Email
	}

	updatedUser, err := app.store.Users.UpdateUser(c.Request.Context(), user, authUser.ID)
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

// UpdatePassword godoc
//
//	@Summary		Update User Password
//	@Description	Update Current User Password
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		updatePasswordRequest	true	"update credentials"
//	@Success		201		{object}	string					"password updated"
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/auth/change-password [put]
//
//	@Security		BearerAuth
func (app *application) updatePassword(c *gin.Context) {

	var payload updatePasswordRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authUser, err := app.getUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(authUser.Password), []byte(payload.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	if payload.NewPassword != payload.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "passwords do not match"})
		return
	}

	if payload.NewPassword == payload.OldPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "use a different password"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user := &models.User{
		Password: string(hashedPassword),
	}

	if err := app.store.Users.UpdatePassword(c.Request.Context(), user, authUser.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}

	c.JSON(http.StatusCreated, "password updated successfully")
}
