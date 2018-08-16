package router

import (
	"api"

	"github.com/labstack/echo"
)

func New() *echo.Echo {
	e := echo.New()

	// create groups
	authGroup := e.Group("/v1.0/auth")
	userGroupFWC := e.Group("/v1.0/fwc/")
	userGroupSWC := e.Group("/v1.0/swc/")

	// set main routes
	api.MainGroup(e)

	// set group routes
	api.UserGroupFWC(userGroupFWC)
	api.UserGroupSWC(userGroupSWC)
	api.AuthGroup(authGroup)

	return e
}
