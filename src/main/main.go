package main

import (
	"database/sql"
	"fmt"
	"router"

	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	fmt.Println("Welcome to the webserver")

	e := router.New()

	var err error
	db, err = sql.Open("postgres", "postgres://airin:password@localhost/bookstore?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You connected to your database.")

	e.Start(":8000")

}
