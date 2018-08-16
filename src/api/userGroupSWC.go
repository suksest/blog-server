package api

import (
	"api/handlers"
	swc "api/middlewares/ratelimiter/slidingwindowcounter"

	"github.com/labstack/echo"
)

func UserGroupSWC(g *echo.Group) {

	// config := fixedwindow.NewConfig("userlimiter", 5, "minute")
	// g.Use(fixedwindow.UserLimiter(config))
	// configTB := tokenbucket.NewConfig("userlimitertb", 3, "minute")
	// g.Use(tokenbucket.Limiter(configTB))

	config := swc.NewConfig("userlimiterswc", 5, "minute")
	g.Use(swc.UserLimiter(config))

	g.GET("/users", handlers.GetAllUser)
	g.GET("/user/:id", handlers.GetUserByID)
	g.DELETE("/user/:id", handlers.DeleteUserByID)
	g.PUT("/user/:id", handlers.UpdateUser)
}
