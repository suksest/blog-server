package api

import (
	"api/handlers"
	swl "api/middlewares/ratelimiter/slidingwindowlog"

	"github.com/labstack/echo"
)

func SwlGroup(g *echo.Group) {

	// config := fixedwindow.NewConfig("userlimiter", 5, "minute")
	// g.Use(fixedwindow.UserLimiter(config))
	// configTB := tokenbucket.NewConfig("userlimitertb", 3, "minute")
	// g.Use(tokenbucket.Limiter(configTB))

	config := swl.NewConfig("userlimiterswl", 5, "minute")
	g.Use(swl.UserLimiter(config))

	g.POST("/publish", handlers.PublishPost)
	g.GET("/posts", handlers.GetAllPost)
	g.GET("/post/:id", handlers.GetPostByID)
	g.GET("/author/:id/posts", handlers.GetPostByAuthorID)
	g.DELETE("/post/:id", handlers.DeletePostByID)
	g.PUT("/post/:id", handlers.UpdatePost)
}
