package api

import (
	"api/handlers"

	"github.com/labstack/echo"
)

func MainGroup(e *echo.Echo) {
	//General
	e.GET("/", handlers.Home)

	//User
	e.POST("/user", handlers.SignupUser)
	e.POST("/user/login", handlers.LoginUser)
	e.GET("/user", handlers.GetAllUser)
	e.GET("/user/:id", handlers.GetUserByID)

	//Post
	e.POST("/publish", handlers.PublishPost)

}
