package api

import (
	"api/handlers"

	"github.com/labstack/echo"
)

func MainGroup(e *echo.Echo) {
	e.GET("/login", handlers.Login)
	e.GET("/", handlers.Home)
	e.GET("/user/:data", handlers.GetUser)

	e.POST("/user", handlers.AddUser)
	e.POST("/message", handlers.AddMessage)
	e.POST("/news", handlers.AddNews)

	e.GET("/books", handlers.BooksIndex)
	// e.GET("/books/show", handlers.booksShow)
	// e.GET("/books/create", handlers.booksCreateForm)
	// e.POST("/books/create/process", handlers.booksCreateProcess)
	// e.GET("/books/update", handlers.booksUpdateForm)
	// e.PUT("/books/update/process", handlers.booksUpdateProcess)
	// e.DELETE("/books/delete/process", handlers.booksDeleteProcess)
}
