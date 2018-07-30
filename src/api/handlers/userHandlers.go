package handlers

import (
	"api/postgres"
	"api/user"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

type PasswordMismatchError struct{}

func (e *PasswordMismatchError) Error() string {
	return "password didn't match"
}

func SignupUser(c echo.Context) error {

	theUser := user.User{}

	defer c.Request().Body.Close()

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Failed reading the request body for addUsers: %s\n", err)
		return c.JSON(http.StatusInternalServerError, "Failed reading the request body")
	}

	err = json.Unmarshal(b, &theUser)
	if err != nil {
		log.Printf("Failed unmarshaling in addUsers: %s\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status":  "FAILED",
			"message": "Failed unmarshaling input",
		})
	}

	res, err := Signup(&theUser)
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
	fmt.Println("Created: ", res)
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "OK",
		"message": "User " + fmt.Sprint(res) + " Sucessfully Created",
	})
}

func Signup(req *user.User) (uint, error) {

	db := postgres.OpenDB()
	defer db.Close()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	newUser := &user.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(passwordHash),
	}

	id, err := user.Create(db, newUser)
	if err != nil {
		return 0, err
	}
	return id, err
}

func LoginUser(c echo.Context) error {

	theUser := user.User{}

	defer c.Request().Body.Close()

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Failed reading the request body for addUsers: %s\n", err)
		return c.JSON(http.StatusInternalServerError, "Failed reading the request body")
	}

	err = json.Unmarshal(b, &theUser)
	if err != nil {
		log.Printf("Failed unmarshaling in addUsers: %s\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status":  "FAILED",
			"message": "Failed unmarshaling input",
		})
	}

	res, err := Login(&theUser)
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
	fmt.Printf("Ok: User '%s' logged in\n", res.Username)
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "OK",
		"message": "User " + res.Username + " logged in",
	})

}

func Login(req *user.User) (*user.User, error) {

	db := postgres.OpenDB()
	defer db.Close()

	user, err := user.FindByEmail(db, req.Email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.PasswordHash))
	if err != nil {
		return nil, &PasswordMismatchError{}
	}
	return user, nil
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
	user := user.User{}

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
