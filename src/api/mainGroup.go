package api

import (
	"api/handlers"

	"github.com/labstack/echo"
)

func MainGroup(e *echo.Echo) {
	//General
	e.GET("/", handlers.Home)

	//User
	e.POST("/v1.0/user", handlers.SignupUser)
	e.POST("/v1.0/user/login", handlers.LoginUser)
	e.GET("/v1.0/user", handlers.GetAllUser)
	e.GET("/v1.0/user/:id", handlers.GetUserByID)
	e.DELETE("/v1.0/user/:id", handlers.DeleteByID)
	e.PUT("/v1.0/user/:id", handlers.UpdateUser)

	//Post
	e.POST("/v1.0/publish", handlers.PublishPost)

}
