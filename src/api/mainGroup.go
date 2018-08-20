package api

import (
	"api/handlers"
	"api/middlewares"

	"github.com/labstack/echo"
)

func MainGroup(e *echo.Echo) {

	// set all middlewares
	middlewares.SetMainMiddlewares(e)

	e.POST("/v1.0/publish", handlers.PublishPost)
	e.GET("/v1.0/posts", handlers.GetAllPost)
	e.GET("/v1.0/post/:id", handlers.GetPostByID)
	e.GET("/v1.0/author/:id/posts", handlers.GetPostByAuthorID)
	e.DELETE("/v1.0/post/:id", handlers.DeletePostByID)
	e.PUT("/v1.0/post/:id", handlers.UpdatePost)

}
