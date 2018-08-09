package router

import (
	"api"
	"api/middlewares"
	limiter "api/middlewares/ratelimiter/fixedwindowcounter"

	"github.com/labstack/echo"
)

func New() *echo.Echo {
	e := echo.New()

	// create groups
	jwtGroup := e.Group("/jwt")

	// set all middlewares
	middlewares.SetMainMiddlewares(e)
	middlewares.SetJwtMiddlewares(jwtGroup)

	config := limiter.NewConfig("userlimiter", 3, "minute")
	e.Use(limiter.Limiter(config))

	// set main routes
	api.MainGroup(e)

	// set group routes
	api.JwtGroup(jwtGroup)

	return e
}
