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
	e.DELETE("/user/:id", handlers.DeleteByID)
	e.PUT("/user/:id", handlers.UpdateUser)

	//Post
	e.POST("/publish", handlers.PublishPost)

}
