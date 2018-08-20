package api

import (
	"api/handlers"
	fwc "api/middlewares/ratelimiter/fixedwindowcounter"

	"github.com/labstack/echo"
)

func FwcGroup(g *echo.Group) {

	config := fwc.NewConfig("userlimiterfwc", 5, "minute")
	g.Use(fwc.UserLimiter(config))

	g.POST("/publish", handlers.PublishPost)
	g.GET("/posts", handlers.GetAllPost)
	g.GET("/post/:id", handlers.GetPostByID)
	g.GET("/author/:id/posts", handlers.GetPostByAuthorID)
	g.DELETE("/post/:id", handlers.DeletePostByID)
	g.PUT("/post/:id", handlers.UpdatePost)
}
