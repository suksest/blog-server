package api

import (
	"api/handlers"

	"github.com/labstack/echo"
)

func MainGroup(e *echo.Echo) {
	e.GET("/login", handlers.LoginJWT)
	e.GET("/", handlers.Home)
	e.GET("/user/:data", handlers.GetUser)

	e.POST("/user", handlers.AddUser)
	e.POST("/user/signup", handlers.SignupUser)
	e.POST("/user/login", handlers.LoginUser)
	e.POST("/message", handlers.AddMessage)
	e.POST("/news", handlers.AddNews)

	e.POST("/publish", handlers.PublishPost)

}
