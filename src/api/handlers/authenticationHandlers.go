package handlers

import (
	"api/user"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"redis"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type JwtClaims struct {
	Hash string `json:"hash"`
	jwt.StandardClaims
}

func createJwtToken(s string) (string, error) {
	claims := JwtClaims{
		s,
		jwt.StandardClaims{
			Id:        "main_user_id",
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	token, err := rawToken.SignedString([]byte("mySecret"))
	if err != nil {
		return "", err
	}

	return token, nil
}

func SignupUser(c echo.Context) error {

	theUser := user.User{}

	defer c.Request().Body.Close()

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Failed reading the request body for sign up: %s\n", err)
		return c.JSON(http.StatusInternalServerError, "Failed reading the request body")
	}

	err = json.Unmarshal(b, &theUser)
	if err != nil {
		log.Printf("Failed unmarshaling in signup user: %s\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status":  "FAILED",
			"message": "Failed unmarshaling input",
		})
	}

	res, err := user.Signup(&theUser)
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
		"message": "User with ID: " + fmt.Sprint(res) + " Sucessfully Created",
	})
}

func LoginUser(c echo.Context) error {

	theUser := user.User{}

	// fmt.Printf("RealIP:" + fmt.Sprint(c.RealIP()) + "\n")
	// fmt.Printf("Time:" + fmt.Sprint(time.Now().Unix()) + "\n")

	defer c.Request().Body.Close()

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Failed reading the request body for Login user: %s\n", err)
		return c.JSON(http.StatusInternalServerError, "Failed reading the request body")
	}

	err = json.Unmarshal(b, &theUser)
	if err != nil {
		log.Printf("Failed unmarshaling in Login user: %s\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status":  "FAILED",
			"message": "Failed unmarshaling input",
		})
	}

	res, err := user.Login(&theUser)
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

	response := c.Response()
	token, err := createJwtToken(res.Username)
	if err != nil {
		log.Println("Error Creating JWT token", err)
		return c.String(http.StatusInternalServerError, "something went wrong")
	}
	response.Header().Set(echo.HeaderAuthorization, token)

	r := redis.RedisConnect()
	defer r.Close()

	t, err := json.Marshal(token)
	if err != nil {
		panic(err)
	}

	// Save JSON blob to Redis
	_, err = r.Do("SET", "user:"+res.Username, t)
	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status":  "OK",
		"message": "User " + res.Username + " logged in",
	})
}
