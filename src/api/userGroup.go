package api

import (
	"api/handlers"
	fixedwindow "api/middlewares/ratelimiter/fixedwindowcounter"

	"github.com/labstack/echo"
)

func UserGroup(g *echo.Group) {

	config := fixedwindow.NewConfig("userlimiter", 5, "minute")
	g.Use(fixedwindow.UserLimiter(config))

	g.GET("/users", handlers.GetAllUser)
	g.GET("/user/:id", handlers.GetUserByID)
	g.DELETE("/user/:id", handlers.DeleteUserByID)
	g.PUT("/user/:id", handlers.UpdateUser)
}
