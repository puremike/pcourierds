package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"slices"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/puremike/pcourierds/internal/store"
)

func (app *application) basicAuthentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.Header("WWW-Authenticate", `Basic realm="restricted"`)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "basic" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is deformed"})
			c.Abort()
			return
		}

		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid base64 encoding"})
			c.Abort()
			return
		}

		username := app.config.basicAuthConfig.username
		password := app.config.basicAuthConfig.password

		decodedStr := string(decoded)

		creds := strings.SplitN(decodedStr, ":", 2)

		if len(creds) != 2 || creds[0] != username || creds[1] != password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (app *application) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is deformed"})
			c.Abort()
			return
		}

		token := strings.TrimSpace(parts[1])

		jwtToken, err := app.jwtAuth.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if jwtToken == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)
		userId, ok := claims["sub"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid sub type"})
			c.Abort()
			return
		}

		user, err := app.store.Users.GetUserById(c.Request.Context(), userId)
		if err != nil {
			if errors.Is(err, store.ErrUserNotFound) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Set("userId", user.ID)
		c.Next()
	}
}

func (app *application) authorizeRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := app.getUserFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		if slices.Contains(allowedRoles, user.Role) {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient role"})
	}
}

func (app *application) getDispatcherAppMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("id")

		dispatchApps, err := app.store.DispatcherApplications.GetApplicationById(c.Request.Context(), id)

		if err != nil {
			if errors.Is(err, store.ErrDispatcherApplicationNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "dispatcher application not found"})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve dispatcher application"})
			c.Abort()
			return
		}

		c.Set("dispatcherApp", dispatchApps)
		c.Next()
	}
}

func (app *application) getDispatcherAppByUserIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		id := c.Param("userID")

		fmt.Println("id: ", id)

		dispatchApp, err := app.store.DispatcherApplications.GetApplicationByUserId(c.Request.Context(), id)

		if err != nil {
			if errors.Is(err, store.ErrDispatcherApplicationNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "dispatcher applicationzsd not found"})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve dispatcher application"})
			c.Abort()
			return
		}

		c.Set("dispatcherAppByUserId", dispatchApp)
		c.Next()
	}
}
