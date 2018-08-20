package router

import (
	"api"

	"github.com/labstack/echo"
	cho "github.com/webx-top/echo"
)

func New() *echo.Echo {
	e := echo.New()

	// create groups
	authGroup := e.Group("/v1.0/auth")
	userGroup := e.Group("/v1.0")

	// set main routes
	api.MainGroup(e)

	// set group routes
	api.UserGroup(userGroup)
	api.AuthGroup(authGroup)

	//bencmarking needs
	swcGroup := e.Group("/v1.0/swc")
	swlGroup := e.Group("/v1.0/swl")
	fwcGroup := e.Group("/v1.0/fwc")
	tbGroup := e.Group("/v1.0/tb")
	tollboothGroup := e.Group("/v1.0/tollbooth")

	g := cho.New()
	webxGroup := g.Group("/v1.0/webx")

	api.SwcGroup(swcGroup)
	api.SwlGroup(swlGroup)
	api.FwcGroup(fwcGroup)
	api.TbGroup(tbGroup)
	api.TollboothGroup(tollboothGroup)
	api.WebxGroup(webxGroup)

	return e
}
