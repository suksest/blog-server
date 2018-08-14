package api

import (
	"api/handlers"

	// fixedwindow "api/middlewares/ratelimiter/fixedwindowcounter"
	swl "api/middlewares/ratelimiter/slidingwindowlog"

	"github.com/labstack/echo"
)

func AuthGroup(g *echo.Group) {
	// configAnon := fixedwindow.NewConfig("anonlimiter", 2, "minute")
	// g.Use(fixedwindow.AnonLimiter(configAnon))
	configAnon := swl.NewConfig("anonlimiterswl", 5, "minute")
	g.Use(swl.AnonLimiter(configAnon))

	g.POST("/user", handlers.SignupUser)
	g.POST("/user/login", handlers.LoginUser)
}
