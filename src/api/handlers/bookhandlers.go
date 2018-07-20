package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

var db *sql.DB

// export fields to templates
// fields changed to uppercase
type Book struct {
	Isbn   string  `json:"Isbn"`
	Title  string  `json:"Title"`
	Author string  `json:"Author"`
	Price  float64 `json:"Price"`
}

func BooksIndex(c echo.Context) error {

	var err error
	db, err = sql.Open("postgres", "postgres://airin:password@localhost/bookstore?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	} else {
		fmt.Println("You connected to your database.")
	}

	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "you need to lets us know if you want json or string data",
		})
	}
	defer rows.Close()

	bks := make([]map[string]string, 0)
	for rows.Next() {
		bk := Book{}
		err := rows.Scan(&bk.Isbn, &bk.Title, &bk.Author, &bk.Price) // order matters
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "you need to lets us know if you want json or string data",
			})
		}
		book := map[string]string{
			"Isbn":   bk.Isbn,
			"Title":  bk.Title,
			"Author": bk.Author,
			"Price":  strconv.FormatFloat(bk.Price, 'f', 6, 64),
		}
		bks = append(bks, book)
	}
	if err = rows.Err(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "you need to lets us know if you want json or string data",
		})
	}

	return c.JSON(http.StatusOK, bks)

}
