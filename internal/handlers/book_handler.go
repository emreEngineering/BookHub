package handlers

import (
	"BookHub/internal/services"
	"encoding/json"
	"net/http"
	"strconv"

	"BookHub/internal/models"
	"BookHub/internal/repositories"
)

type BookHandler struct {
	bookService services.BookServices
}

func NewBookHandler(bookRepo repositories.BookRepository) *BookHandler {
	return &BookHandler{
		bookService: services.NewBookService(bookRepo),
	}
}

func (h *BookHandler) BooksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		var book models.Book

		err := json.NewDecoder(r.Body).Decode(&book)
		if err != nil {
			http.Error(w, "Geçersiz JSON", http.StatusBadRequest)
			return
		}

		createdBook, err := h.bookService.CreateBook(book)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdBook)
		return
	}

	if r.Method == http.MethodPut {
		idParam := r.URL.Query().Get("id")
		if idParam == "" {
			http.Error(w, "Kitap ID zorunludur", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "Geçersiz kitap ID", http.StatusBadRequest)
			return
		}

		var book models.Book

		err = json.NewDecoder(r.Body).Decode(&book)
		if err != nil {
			http.Error(w, "Geçersiz JSON", http.StatusBadRequest)
			return
		}

		updatedBook, err := h.bookService.UpdateBook(id, book)
		if err != nil {
			if err.Error() == "Kitap bulunamadı" {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		json.NewEncoder(w).Encode(updatedBook)
		return
	}

	if r.Method == http.MethodDelete {
		idParam := r.URL.Query().Get("id")
		if idParam == "" {
			http.Error(w, "Kitap ID zorunludur", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "Geçersiz kitap ID", http.StatusBadRequest)
			return
		}

		err = h.bookService.DeleteBook(id)
		if err != nil {
			http.Error(w, "Kitap bulunamadı", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Bu endpoint sadece GET, POST, PUT ve DELETE destekler", http.StatusMethodNotAllowed)
		return
	}

	idParam := r.URL.Query().Get("id")

	if idParam != "" {
		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "Geçersiz kitap ID", http.StatusBadRequest)
			return
		}

		book, err := h.bookService.GetBookByID(id)
		if err != nil {
			http.Error(w, "Kitap bulunamadı", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(book)
		return
	}

	books, err := h.bookService.GetAllBooks()
	if err != nil {
		http.Error(w, "Kitaplar alınamadı", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(books)
}
