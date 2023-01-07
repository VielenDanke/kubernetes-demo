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

func parseFileEnvConfig() (config, error) {
	var cfg config

	file, err := os.Open("/tmp/config.json")

	if err != nil {
		log.Println("File is not found, parse env")

		appPort := os.Getenv("APPLICATION_PORT")
		appName := os.Getenv("APPLICATION_NAME")

		if len(appPort) == 0 || len(appName) == 0 {
			return cfg, fmt.Errorf("env config are not exists")
		}
		cfg.Port = appPort
		cfg.Name = appName

		return cfg, nil
	} else {
		log.Println("File is found - parse file")

		if err = json.NewDecoder(file).Decode(&cfg); err != nil {
			return cfg, err
		}
	}
	return cfg, nil
}

func setupDefaultConfig(cfg *config) {
	log.Println("File and Env config are not found. Setup default")
	if len(cfg.Port) == 0 {
		cfg.Port = "8080"
	}
	if len(cfg.Name) == 0 {
		cfg.Name = "books"
	}
	return
}

func main() {
	cfg, err := parseFileEnvConfig()

	if err != nil {
		log.Printf("Config are not set up. The reason: %v. Using default instead.", err)
		setupDefaultConfig(&cfg)
	}

	log.Printf("Application %s is starting\n", cfg.Name)

	log.Printf("DB credentials: username: %s, password: %s\n", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"))

	router := gin.Default()

	router.GET("/books", findAll)

	router.GET("/books/:id", findById)

	router.POST("/books", save)

	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	log.Fatalln(router.Run(fmt.Sprintf(":%s", cfg.Port)))
}
