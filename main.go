package main

import (
	"encoding/json"
	"log"
	"net/http"
)

var books []Book

type Book struct {
	Name string `json:"name"`
}

func main() {
	srv := http.NewServeMux()

	srv.HandleFunc("/health", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})
	srv.HandleFunc("/books", func(rw http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var book Book
			if decodeErr := json.NewDecoder(r.Body).Decode(&book); decodeErr != nil {
				errorMessage := decodeErr.Error()
				log.Printf("ERROR: cannot decode json, message: %s\n", errorMessage)
				rw.Write([]byte(errorMessage))
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
			books = append(books, book)
			log.Printf("INFO: book with name %s successfully saved\n", book.Name)
			rw.WriteHeader(http.StatusCreated)
		case http.MethodGet:
			if encodeErr := json.NewEncoder(rw).Encode(&books); encodeErr != nil {
				errorMessage := encodeErr.Error()
				log.Printf("ERROR: cannot encode data to json, message: %s\n", errorMessage)
				rw.Write([]byte(errorMessage))
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
		default:
			rw.WriteHeader(http.StatusNotFound)
		}
	})
	log.Printf("INFO: Server started on port: %s\n", "9090")

	log.Fatalln(http.ListenAndServe(":9090", srv))
}
