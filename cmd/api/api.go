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
		api.GET("/health", app.health)
		api.GET("swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	auth := api.Group("/auth")
	{
		auth.POST("/signup", app.CreateUser)
		// auth.POST("/login", app.Login)
	}
	return g
}
