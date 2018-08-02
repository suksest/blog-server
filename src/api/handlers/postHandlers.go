package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"redis"
	"strconv"
	"time"

	"api/post"
	"api/tag"
	"api/user"
	"postgres"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
)

type Postingan struct {
	Title    string
	Body     string
	AuthorID uint
	Tags     []string
}

type Response struct {
	PostId uint
}

func PublishPost(c echo.Context) error {
	req := c.Request()
	token := req.Header.Get("Authorization")

	if redis.Find(token) {
		postingan := Postingan{}

		defer c.Request().Body.Close()

		b, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			log.Printf("Failed reading the request body for publish post: %s\n", err)
			return c.JSON(http.StatusInternalServerError, "Failed reading the request body")
		}

		err = json.Unmarshal(b, &postingan)
		if err != nil {
			log.Printf("Failed unmarshaling in publish post: %s\n", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"status":  "FAILED",
				"message": "Failed unmarshaling input",
			})
		}

		// #dummy
		// postingan := Postingan{
		// 	AuthorID: 1,
		// 	Body:     "Go golang rocks! ",
		// 	Title:    "My gomidway post2",
		// 	Tags:     []string{"intro", "golang"},
		// }

		res, err := NewPost(&postingan)
		fmt.Println(res)
		if err != nil {
			if _, ok := err.(*post.TitleDuplicateError); ok {
				// fmt.Println("Bad Request: ", err.Error())
				return c.JSON(http.StatusBadRequest, map[string]string{
					"status":  "FAILED",
					"message": "Title already exist",
				})
			}

		}
		return c.JSON(http.StatusOK, map[string]string{
			"status": "OK",
			"value":  postingan.Title,
		})
	} else {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status":  "FORBIDDEN",
			"message": "Invalid Token",
		})
	}

}

func NewPost(postingan *Postingan) (*Response, error) {

	db := postgres.OpenDB()
	defer db.Close()

	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	newPost := &post.Post{
		AuthorID:    postingan.AuthorID,
		Title:       postingan.Title,
		Body:        postingan.Body,
		PublishedAt: time.Now().UTC(),
	}
	_, err := post.Create(tx, newPost)
	if err != nil {
		return nil, err
	}
	for _, tagName := range postingan.Tags {
		t, err := tag.CreateIfNotExists(tx, tagName)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		err = post.AddTag(tx, newPost, t)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	res := tx.Commit()
	if res.Error != nil {
		return nil, res.Error
	}
	return &Response{PostId: newPost.ID}, nil
}

func GetAllPost(c echo.Context) error {

	req := c.Request()
	token := req.Header.Get("Authorization")

	if redis.Find(token) {
		posts, err := post.FindAll()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Failed reading the request body")
		}

		return c.JSON(http.StatusOK, posts)
	} else {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status":  "FORBIDDEN",
			"message": "Invalid Token",
		})
	}

}

func GetPostByID(c echo.Context) error {

	req := c.Request()
	token := req.Header.Get("Authorization")

	if redis.Find(token) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}

		post, err := post.FindByID(uint(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, post)
	} else {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status":  "FORBIDDEN",
			"message": "Invalid Token",
		})
	}

}

func GetPostByAuthorID(c echo.Context) error {

	req := c.Request()
	token := req.Header.Get("Authorization")

	if redis.Find(token) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		fmt.Printf(fmt.Sprint(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}

		posts, err := post.FindByAuthorID(uint(id))
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, posts)
	} else {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status":  "FORBIDDEN",
			"message": "Invalid Token",
		})
	}
}

func DeletePostByID(c echo.Context) error {

	req := c.Request()
	token := req.Header.Get("Authorization")

	if redis.Find(token) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}

		post.Delete(uint(id))

		return c.JSON(http.StatusOK, map[string]string{
			"status":  "OK",
			"message": "Post with ID:" + fmt.Sprint(id) + " sucesfully deleted",
		})
	} else {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status":  "FORBIDDEN",
			"message": "Invalid Token",
		})
	}
}

func UpdatePost(c echo.Context) error {
	req := c.Request()
	token := req.Header.Get("Authorization")

	if redis.Find(token) {
		thePost := post.Post{}

		defer c.Request().Body.Close()

		b, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			log.Printf("Failed reading the request body for update post: %s\n", err)
			return c.JSON(http.StatusInternalServerError, "Failed reading the request body")
		}

		err = json.Unmarshal(b, &thePost)
		if err != nil {
			log.Printf("Failed unmarshaling in update post: %s\n", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"status":  "FAILED",
				"message": "Failed unmarshaling input",
			})
		}

		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}

		res, err := post.Update(&thePost, uint(id))
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
			"message": "Post with ID: " + fmt.Sprint(res.ID) + " Sucessfully Updated",
		})

	} else {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status":  "FORBIDDEN",
			"message": "Invalid Token",
		})
	}
}
