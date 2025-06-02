package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/puremike/pcourierds/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (app *application) routes() http.Handler {

	g := gin.Default()

	docs.SwaggerInfo.BasePath = "/api/v1"
	api := g.Group("/api/v1")
	{
		api.GET("/health", app.basicAuthentication(), app.health)
		api.GET("swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	users := api.Group("/auth")
	{
		users.POST("/signup", app.createUser)
		users.POST("/login", app.login)
	}

	authGroup := api.Group("/")
	authGroup.Use(app.authMiddleware())
	{
		authGroup.GET("/auth/me", app.userProfile)
		authGroup.PATCH("/auth/update-profile", app.updateProfile)
		authGroup.PUT("/auth/change-password", app.updatePassword)

		authGroup.POST("/dispatchers/apply", app.dispatcherApply)
		authGroup.GET("/admin/dispatcher-applications", app.authorizeRoles("admin"), app.getAllApplications)
		authGroup.GET("/admin/dispatcher-applications/:id", app.authorizeRoles("admin"), app.getDispatcherAppMiddleware(), app.getDispatcherApplicationById)
		authGroup.PATCH("/admin/approve-dispatcher/:userID", app.authorizeRoles("admin"), app.getDispatcherAppByUserIdMiddleware(), app.approveDenyApplication)

		authGroup.GET("/admin/user/:id", app.authorizeRoles("user", "admin"), app.getUserById)
		authGroup.GET("/admin/users", app.authorizeRoles("admin"), app.getUsers)
		authGroup.POST("/admin/user", app.authorizeRoles("admin"), app.adminCreateUser)
		authGroup.PATCH("/admin/user/:id", app.authorizeRoles("admin"), app.adminUpdateProfile)
		authGroup.DELETE("/admin/user/:id", app.authorizeRoles("admin"), app.adminDeleteUser)
	}

	return g
}
