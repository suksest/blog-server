package handlers

import (
	"api/user"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"rateLimiter/fixedWindowCounter"
	"redis"
	"strconv"

	"github.com/labstack/echo"
)

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
	token, err := createJwtToken()
	if err != nil {
		log.Println("Error Creating JWT token", err)
		return c.String(http.StatusInternalServerError, "something went wrong")
	}
	response.Header().Set("Token", token)

	r := redis.RedisConnect()
	defer r.Close()

	t, err := json.Marshal(token)
	if err != nil {
		panic(err)
	}

	// Save JSON blob to Redis
	reply, err := r.Do("SET", "user:"+res.Username, t)
	if err != nil {
		panic(err)
	}

	if fixedWindowCounter.UserLimiter(res.Username) {
		fmt.Printf("user:" + res.Username + " " + fmt.Sprint(reply))
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "OK",
			"message": "User " + res.Username + " logged in",
		})
	} else {
		fmt.Printf("user:" + res.Username + " " + fmt.Sprint(reply))
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "FORBIDDEN",
			"message": "Request Limit exceeded",
		})
	}
}

func GetAllUser(c echo.Context) error {

	req := c.Request()
	token := req.Header.Get("Authorization")

	if redis.Find(token) {
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

	if redis.Find(token) {
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

	if redis.Find(token) {
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

	if redis.Find(token) {
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
