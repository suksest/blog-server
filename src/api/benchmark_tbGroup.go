package api

import (
	"api/handlers"
	tb "api/middlewares/ratelimiter/tokenbucket"

	"github.com/labstack/echo"
)

func TbGroup(g *echo.Group) {

	config := tb.NewConfig("userlimitertb", 5, "minute")
	g.Use(tb.UserLimiter(config))

	g.POST("/publish", handlers.PublishPost)
	g.GET("/posts", handlers.GetAllPost)
	g.GET("/post/:id", handlers.GetPostByID)
	g.GET("/author/:id/posts", handlers.GetPostByAuthorID)
	g.DELETE("/post/:id", handlers.DeletePostByID)
	g.PUT("/post/:id", handlers.UpdatePost)
}
