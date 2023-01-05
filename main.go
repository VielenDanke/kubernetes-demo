package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math"
	"net/http"
	"strconv"
)

var books = []book{
	{ID: 1, Name: "first"},
	{ID: 2, Name: "second"},
	{ID: 3, Name: "third"},
}

type book struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func findAll(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func findById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	for _, book := range books {
		if book.ID == id {
			c.IndentedJSON(http.StatusOK, book)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Book with ID %d not found", id)})
}

func save(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	newBook.ID = generateNextID()
	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, gin.H{"book_id": newBook.ID})
}

func generateNextID() int {
	maxID := 0

	for _, v := range books {
		maxID = int(math.Max(float64(v.ID), float64(maxID)))
	}
	return maxID + 1
}

func main() {
	router := gin.Default()

	router.GET("/books", findAll)

	router.GET("/books/:id", findById)

	router.POST("/books", save)

	log.Fatalln(router.Run("localhost:8080"))
}
