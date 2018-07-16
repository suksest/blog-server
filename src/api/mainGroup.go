package api

import (
	"api/handlers"

	"github.com/labstack/echo"
)

func MainGroup(e *echo.Echo) {
	e.GET("/login", handlers.Login)
	e.GET("/home", handlers.Home)
	e.GET("/user/:data", handlers.GetUser)

	e.POST("/user", handlers.AddUser)
	e.POST("/message", handlers.AddMessage)
	e.POST("/news", handlers.AddNews)
}
