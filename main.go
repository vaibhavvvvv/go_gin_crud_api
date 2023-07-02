package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "In search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "war and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

func getBooks(c *gin.Context) { //gin.context stoes all the info related to the specific request
	c.IndentedJSON(http.StatusOK, books)
}

func createBook(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("book not found")
}

func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
	}
	c.IndentedJSON(http.StatusOK, book)

}

func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing id in query"})
		return
	}

	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return

	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "book not available"})
		return
	}

	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}

func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing id in query"})
		return
	}

	book, err := getBookById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return

	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

func deleteBook(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing id in query"})
		return
	}

	for i, b := range books {
		if id == b.ID {
			books = append(books[:i], books[i+1:]...)
			c.IndentedJSON(http.StatusOK, books)
			return

		}
	}

	c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available"})

}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)          //curl localhost:8080/books
	router.POST("/books", createBook)       //curl localhost:8080/books --include --header "Content-Type: application/json" -d @body.json --request "POST"
	router.GET("/books/:id", bookById)      // curl localhost:8080/books/3
	router.PATCH("/checkout", checkoutBook) //curl localhost:8080/checkout?id=2 --request "PATCH"
	router.PATCH("/return", returnBook)     //curl localhost:8080/return?id=2 --request "PATCH"
	router.DELETE("/delete", deleteBook)

	router.Run("localhost:8080")
}
