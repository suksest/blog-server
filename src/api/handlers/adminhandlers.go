package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

func MainAdmin(c echo.Context) error {
	return c.String(http.StatusOK, "This is the admin main page")
}
