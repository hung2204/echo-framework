package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/warrant-dev/warrant-go/v4"
)

// Book represents a book entity.
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books = make(map[int]Book)
var currentID = 0

func createBook(c echo.Context) error {
	var newBook Book
	if err := c.Bind(&newBook); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	currentID++
	newBook.ID = currentID
	books[newBook.ID] = newBook

	log.Printf("createBook: %v", newBook)
	return c.JSON(http.StatusCreated, newBook)
}

func getBooks(c echo.Context) error {
	log.Printf("getBooks: %v", books)
	return c.JSON(http.StatusOK, books)
}

func getBook(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid book ID")
	}

	book, ok := books[id]
	if !ok {
		return c.JSON(http.StatusNotFound, "Book not found")
	}

	log.Printf("getBook by id: %v", book)
	return c.JSON(http.StatusOK, book)
}

func updateBook(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid book ID")
	}

	var updatedBook Book
	if err := c.Bind(&updatedBook); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	_, ok := books[id]
	if !ok {
		return c.JSON(http.StatusNotFound, "Book not found")
	}

	updatedBook.ID = id
	books[id] = updatedBook

	log.Printf("updateBook: %v", updatedBook)
	return c.JSON(http.StatusOK, updatedBook)
}

func deleteBook(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid book ID")
	}

	_, ok := books[id]
	if !ok {
		return c.JSON(http.StatusNotFound, "Book not found")
	}

	delete(books, id)

	log.Printf("deleteBook success")
	return c.NoContent(http.StatusNoContent)
}

func main() {
	e := echo.New()

	warrant.ApiKey = "api_test_f5dsKVeYnVSLHGje44zAygqgqXiLJBICbFzCiAg1E="
	warrant.ApiEndpoint = "http://localhost:8080"
	warrant.AuthorizeEndpoint = "http://localhost:8080"

	// Create some initial books
	books[1] = Book{ID: 1, Title: "Book 1", Author: "Author 1"}
	books[2] = Book{ID: 2, Title: "Book 2", Author: "Author 2"}

	// Routes for CRUD operations
	e.POST("/books", createBook)
	e.GET("/books", getBooks)
	e.GET("/books/:id", getBook)
	e.PUT("/books/:id", updateBook)
	e.DELETE("/books/:id", deleteBook)

	e.Logger.Fatal(e.Start(":8080"))
}
