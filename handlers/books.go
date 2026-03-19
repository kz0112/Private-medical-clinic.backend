package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"Private-medical-clinic.backend/models"
	"Private-medical-clinic.backend/storage"
)

func BooksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case http.MethodGet:
		books := storage.Books

		// Filters
		author := strings.ToLower(r.URL.Query().Get("author"))
		category := strings.ToLower(r.URL.Query().Get("category"))
		title := strings.ToLower(r.URL.Query().Get("title"))

		var filtered []models.Book

		for _, book := range books {
			if author != "" && !strings.Contains(strings.ToLower(book.Author), author) {
				continue
			}
			if category != "" && !strings.Contains(strings.ToLower(book.Category), category) {
				continue
			}
			if title != "" && !strings.Contains(strings.ToLower(book.Title), title) {
				continue
			}
			filtered = append(filtered, book)
		}

		// Pagination
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

		if page <= 0 {
			page = 1
		}
		if limit <= 0 {
			limit = 5
		}

		start := (page - 1) * limit
		end := start + limit

		if start > len(filtered) {
			start = len(filtered)
		}
		if end > len(filtered) {
			end = len(filtered)
		}

		json.NewEncoder(w).Encode(filtered[start:end])

	case http.MethodPost:
		var newBook models.Book

		err := json.NewDecoder(r.Body).Decode(&newBook)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if newBook.Title == "" || newBook.Author == "" {
			http.Error(w, "Title and Author required", http.StatusBadRequest)
			return
		}

		newBook.ID = len(storage.Books) + 1
		storage.Books = append(storage.Books, newBook)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newBook)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func BookByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := strings.TrimPrefix(r.URL.Path, "/books/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for i, book := range storage.Books {
		if book.ID == id {

			switch r.Method {

			case http.MethodGet:
				json.NewEncoder(w).Encode(book)

			case http.MethodPut:
				var updatedBook models.Book

				err := json.NewDecoder(r.Body).Decode(&updatedBook)
				if err != nil {
					http.Error(w, "Invalid JSON", http.StatusBadRequest)
					return
				}

				updatedBook.ID = id
				storage.Books[i] = updatedBook

				json.NewEncoder(w).Encode(updatedBook)

			case http.MethodDelete:
				storage.Books = append(storage.Books[:i], storage.Books[i+1:]...)
				json.NewEncoder(w).Encode(map[string]string{
					"message": "Deleted successfully",
				})

			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}

			return
		}
	}

	http.Error(w, "Book not found", http.StatusNotFound)
}
