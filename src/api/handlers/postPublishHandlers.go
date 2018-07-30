package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"api/post"
	"api/postgres"
	"api/tag"

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

	postingan := Postingan{}

	defer c.Request().Body.Close()

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Failed reading the request body for addUsers: %s\n", err)
		return c.JSON(http.StatusInternalServerError, "Failed reading the request body")
	}

	err = json.Unmarshal(b, &postingan)
	if err != nil {
		log.Printf("Failed unmarshaling in addUsers: %s\n", err)
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
				"message": "Input Not Valid",
			})
		}

	}
	// return c.JSON(http.StatusOK, &postingan)
	// fmt.Println("You connected to your databasea.")
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "OK",
		"message": postingan.Title,
	})
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
