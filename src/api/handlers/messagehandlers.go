package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

type Message struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func AddMessage(c echo.Context) error {
	message := Message{}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&message)
	if err != nil {
		log.Printf("Failed processing addMessage request: %s\n", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	log.Printf("this is your message: %#v", message)
	return c.String(http.StatusOK, "we got your message!")
}
