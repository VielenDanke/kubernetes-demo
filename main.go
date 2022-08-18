package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

type Book struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func initDBConnection(connectionUrl string) (*sql.DB, error) {
	conn, openErr := sql.Open("postgres", connectionUrl)
	if openErr != nil {
		return nil, openErr
	}
	return conn, conn.Ping()
}

func addTable(conn *sql.DB) error {
	_, execErr := conn.Exec("create table if not exists books (id bigserial primary key, name varchar (256))")
	return execErr
}

func main() {
	srv := http.NewServeMux()

	/*
		Deploy with secrets in deployment (image version 2)

		postgresConnectionUrlBytes, _ := os.ReadFile("/tmp/postgres")

		connection, connErr := initDBConnection(string(postgresConnectionUrlBytes))
	*/

	connection, connErr := initDBConnection("host=postgres port=5432 dbname=books user=user password=password sslmode=disable")

	if connErr != nil {
		log.Fatalln(connErr)
	}
	if migrErr := addTable(connection); migrErr != nil {
		log.Fatalln(migrErr)
	}
	srv.HandleFunc("/health", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})
	srv.HandleFunc("/books", func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
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
			tx, _ := connection.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
			_, scErr := tx.ExecContext(ctx, "insert into books (name) values($1)", book.Name)
			if scErr != nil {
				rw.Write([]byte(scErr.Error()))
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			tx.Commit()
			log.Printf("INFO: book with name %s successfully saved\n", book.Name)
			rw.WriteHeader(http.StatusCreated)
		case http.MethodGet:
			books := make([]Book, 0)

			tx, _ := connection.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})

			result, resultErr := tx.QueryContext(ctx, "select b.id, b.name from books b")

			if resultErr != nil {
				rw.Write([]byte(resultErr.Error()))
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			for result.Next() {
				var book Book
				if scErr := result.Scan(&book.ID, &book.Name); scErr != nil {
					rw.Write([]byte(scErr.Error()))
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
				books = append(books, book)
			}
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
	applicationPort, applicationName := os.Getenv("APPLICATION_PORT"), os.Getenv("APPLICATION_NAME")

	log.Printf("INFO: Application %s started on port: %s\n", applicationName, applicationPort)

	log.Fatalln(http.ListenAndServe(applicationPort, srv))
}
