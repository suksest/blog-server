package router

import (
	"api"

	"github.com/labstack/echo"
)

func New() *echo.Echo {
	e := echo.New()

	// create groups
	userGroup := e.Group("/v1.0/")
	authGroup := e.Group("/v1.0/auth")

	// set main routes
	api.MainGroup(e)

	// set group routes
	api.UserGroup(userGroup)
	api.AuthGroup(authGroup)

	return e
}
