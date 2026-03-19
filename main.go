package main

import (
	"log"
	"net/http"

	"Private-medical-clinic.backend/handlers"

	_ "Private-medical-clinic.backend/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	http.HandleFunc("/books", handlers.BooksHandler)
	http.HandleFunc("/books/", handlers.BookByIDHandler)

	http.Handle("/swagger/", httpSwagger.WrapHandler)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
