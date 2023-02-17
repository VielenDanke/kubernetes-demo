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

func deleteById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	for idx, book := range books {
		if book.ID == id {
			books = append(books[:idx], books[(idx+1):]...)
			c.Status(http.StatusOK)
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
		return cfg, err
	}
	log.Println("File is found - parse file")

	if err = json.NewDecoder(file).Decode(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func parseEnvConfig() (config, error) {
	var cfg config

	log.Println("File is not found, parse env")

	appPort := os.Getenv("APPLICATION_PORT")
	appName := os.Getenv("APPLICATION_NAME")

	if len(appPort) == 0 || len(appName) == 0 {
		return cfg, fmt.Errorf("env configuration doesn't exists")
	}
	cfg.Port = appPort
	cfg.Name = appName

	return cfg, nil
}

func setupDefaultConfig(cfg *config) {
	if len(cfg.Port) == 0 {
		cfg.Port = "8080"
	}
	if len(cfg.Name) == 0 {
		cfg.Name = "books"
	}
	return
}

func main() {
	var cfg config
	var err error
	cfg, err = parseFileEnvConfig()
	cfg, err = parseEnvConfig()

	if err != nil {
		log.Println("No Env and File config found. Setup default.")
		setupDefaultConfig(&cfg)
	}
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")

	if len(username) == 0 || len(password) == 0 {
		log.Println("Credentials are not exists")
	} else {
		log.Printf("DB credentials: username: %s, password: %s", username, password)
	}
	router := gin.New()

	router.GET("/books", findAll)

	router.GET("/books/:id", findById)

	router.POST("/books", save)

	router.DELETE("/books/:id", deleteById)

	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	log.Printf("Application %s is starting\n", cfg.Name)

	log.Fatalln(router.Run(fmt.Sprintf(":%s", cfg.Port)))
}
