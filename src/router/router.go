package router

import (
	"api"
	"api/middlewares"

	"github.com/labstack/echo"
)

func New() *echo.Echo {
	e := echo.New()

	// create groups
	jwtGroup := e.Group("/jwt")

	// set all middlewares
	middlewares.SetMainMiddlewares(e)
	middlewares.SetJwtMiddlewares(jwtGroup)

	// set main routes
	api.MainGroup(e)

	// set group routes
	api.JwtGroup(jwtGroup)

	return e
}
