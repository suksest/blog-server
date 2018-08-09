package router

import (
	"api"
	"api/middlewares"
	tokenbucket "api/middlewares/ratelimiter/tokenbucket"

	"github.com/labstack/echo"
)

func New() *echo.Echo {
	e := echo.New()

	// create groups
	jwtGroup := e.Group("/jwt")

	// set all middlewares
	middlewares.SetMainMiddlewares(e)
	middlewares.SetJwtMiddlewares(jwtGroup)

	// config := fixedwindow.NewConfig("userlimiter", 3, "minute")
	configTB := tokenbucket.NewConfig("userlimiter", 3, "minute")
	// e.Use(fixedwindow.Limiter(config))
	e.Use(tokenbucket.Limiter(configTB))

	// set main routes
	api.MainGroup(e)

	// set group routes
	api.JwtGroup(jwtGroup)

	return e
}
