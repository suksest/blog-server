package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

type User struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func GetUser(c echo.Context) error {
	userName := c.QueryParam("name")
	userType := c.QueryParam("type")

	dataType := c.Param("data")

	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("The name is: %s\nand his type is: %s\n", userName, userType))
	}

	if dataType == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"name": userName,
			"type": userType,
		})
	}

	return c.JSON(http.StatusBadRequest, map[string]string{
		"error": "you need to lets us know if you want json or string data",
	})
}

func AddUser(c echo.Context) error {
	user := User{}

	defer c.Request().Body.Close()

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Failed reading the request body for addUsers: %s\n", err)
		return c.String(http.StatusInternalServerError, "")
	}

	err = json.Unmarshal(b, &user)
	if err != nil {
		log.Printf("Failed unmarshaling in addUsers: %s\n", err)
		return c.String(http.StatusInternalServerError, "")
	}

	log.Printf("this is your user: %#v\n", user)
	return c.String(http.StatusOK, "we got your user!")
}
