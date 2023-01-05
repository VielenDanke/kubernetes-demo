package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
)

var books = []book{
	{ID: 1, Name: "first"},
	{ID: 2, Name: "second"},
	{ID: 3, Name: "third"},
}

type config struct {
	Port string `json:"application_port"`
	Name string `json:"application_name"`
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

func parseConfig() (cfg config, err error) {
	file, err := os.Open("/tmp/config.json")

	if err != nil {
		log.Println("File is not found, parse env")

		cfg.Port = os.Getenv("APPLICATION_PORT")
		cfg.Name = os.Getenv("APPLICATION_NAME")
	} else {
		if err = json.NewDecoder(file).Decode(&cfg); err != nil {
			return
		}
	}
	err = validateConfig(cfg)
	return
}

func validateConfig(cfg config) (err error) {
	if len(cfg.Port) == 0 || len(cfg.Name) == 0 {
		err = fmt.Errorf("config is not found")
	}
	return
}

func main() {
	cfg, err := parseConfig()

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Application %s is starting\n", cfg.Name)

	log.Printf("DB credentials: username: %s, password: %s\n", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"))

	router := gin.Default()

	router.GET("/books", findAll)

	router.GET("/books/:id", findById)

	router.POST("/books", save)

	log.Fatalln(router.Run("localhost:%s", cfg.Port))
}
