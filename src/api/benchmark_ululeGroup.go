package api

import (
	"api/handlers"
	"log"

	ulule "api/middlewares/ratelimiter/ulule"

	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"github.com/ulule/limiter"
	sredis "github.com/ulule/limiter/drivers/store/redis"
)

func UluleGroup(g *echo.Group) {

	// Define a limit rate to 4 requests per hour.
	rate, err := limiter.NewRateFromFormatted("5-M")
	if err != nil {
		log.Fatal(err)
		return
	}

	// Create a redis client.
	option, err := redis.ParseURL("redis://localhost:6379/0")
	if err != nil {
		log.Fatal(err)
		return
	}
	client := redis.NewClient(option)

	// Create a store with the redis client.
	store, err := sredis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix:   "userlimiterul",
		MaxRetry: 3,
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	// Create a new middleware with the limiter instance.
	middleware := ulule.NewMiddleware(limiter.New(store, rate))

	g.Use(middleware)

	g.POST("/publish", handlers.PublishPost)
	g.GET("/posts", handlers.GetAllPost)
	g.GET("/post/:id", handlers.GetPostByID)
	g.GET("/author/:id/posts", handlers.GetPostByAuthorID)
	g.DELETE("/post/:id", handlers.DeletePostByID)
	g.PUT("/post/:id", handlers.UpdatePost)
}
