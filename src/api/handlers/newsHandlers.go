package handlers

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
)

type News struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func AddNews(c echo.Context) error {
	news := News{}

	err := c.Bind(&news)
	if err != nil {
		log.Printf("Failed processing addNews request: %s\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	log.Printf("this is your news: %#v", news)
	return c.String(http.StatusOK, "we got your news!")
}
