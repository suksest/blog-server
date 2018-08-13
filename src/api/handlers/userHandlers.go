package handlers

import (
	"api/user"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"redis"
	"strconv"

	"github.com/labstack/echo"
)

func GetAllUser(c echo.Context) error {

	req := c.Request()
	token := req.Header.Get("Authorization")

	if redis.Find(token) != "" {
		users, err := user.FindAll()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Failed reading the request body")
		}

		return c.JSON(http.StatusOK, users)
	} else {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status":  "FORBIDDEN",
			"message": "Invalid Token",
		})
	}

}

func GetUserByID(c echo.Context) error {

	req := c.Request()
	token := req.Header.Get("Authorization")

	if redis.Find(token) != "" {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}

		user, err := user.FindByID(uint(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, user)
	} else {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status":  "FORBIDDEN",
			"message": "Invalid Token",
		})
	}
}

func DeleteUserByID(c echo.Context) error {

	req := c.Request()
	token := req.Header.Get("Authorization")

	if redis.Find(token) != "" {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}

		user.Delete(uint(id))

		return c.JSON(http.StatusOK, map[string]string{
			"status":  "OK",
			"message": "User with ID:" + fmt.Sprint(id) + " sucesfully deleted",
		})
	} else {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status":  "FORBIDDEN",
			"message": "Invalid Token",
		})
	}
}

func UpdateUser(c echo.Context) error {
	req := c.Request()
	token := req.Header.Get("Authorization")

	if redis.Find(token) != "" {
		theUser := user.User{}

		defer c.Request().Body.Close()

		b, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			log.Printf("Failed reading the request body for update user: %s\n", err)
			return c.JSON(http.StatusInternalServerError, "Failed reading the request body")
		}

		err = json.Unmarshal(b, &theUser)
		if err != nil {
			log.Printf("Failed unmarshaling in update user: %s\n", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"status":  "FAILED",
				"message": "Failed unmarshaling input",
			})
		}

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}

		res, err := user.Update(&theUser, uint(id))
		if err != nil {
			switch err.(type) {
			case *user.UsernameDuplicateError:
				fmt.Println("Bad Request: ", err.Error())
				return c.JSON(http.StatusBadRequest, map[string]string{
					"status":  "FAILED",
					"message": "Bad Request",
				})
			case *user.EmailDuplicateError:
				fmt.Println("Bad Request: ", err.Error())
				return c.JSON(http.StatusBadRequest, map[string]string{
					"status":  "FAILED",
					"message": "Bad Request",
				})
			default:
				fmt.Println("Internal Server Error: ", err.Error())
				return c.JSON(http.StatusInternalServerError, map[string]string{
					"status":  "FAILED",
					"message": "Internal Server Error",
				})
			}
		}
		fmt.Println("Updated: ", res.ID)
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "OK",
			"message": "User with ID: " + fmt.Sprint(res.ID) + " Sucessfully Updated",
		})
	} else {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status":  "FORBIDDEN",
			"message": "Invalid Token",
		})
	}
}
