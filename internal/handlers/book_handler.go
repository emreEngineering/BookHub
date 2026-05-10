package handlers

import (
	"BookHub/internal/models"
	"BookHub/internal/repositories"
	"encoding/json"
	"net/http"
	"strconv"
)

var BookRepo repositories.BookRepository

type BookHandler struct{}

func NewBookHandler() *BookHandler {
	return &BookHandler{}
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

		if book.Title == "" {
			http.Error(w, "Kitap adı boş olamaz", http.StatusBadRequest)
			return
		}

		if book.Author == "" {
			http.Error(w, "Yazar adı boş olamaz", http.StatusBadRequest)
			return
		}

		createdBook, err := BookRepo.Create(book)
		if err != nil {
			http.Error(w, "Kitap oluşturulamadı", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdBook)
		return
	}

	if r.Method == http.MethodPut {
		idParam := r.URL.Query().Get("id")
		if idParam == "" {
			http.Error(w, "Kitap ID zorunludur.", http.StatusBadRequest)
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

		if book.Title == "" {
			http.Error(w, "Kitap adı boş olamaz", http.StatusBadRequest)
			return
		}

		if book.Author == "" {
			http.Error(w, "Yazar adı boş olaamaz", http.StatusBadRequest)
			return
		}

		updatedBook, err := BookRepo.Update(id, book)
		if err != nil {
			http.Error(w, "Kitap bulunamadı", http.StatusNotFound)
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

		err = BookRepo.Delete(id)
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

		book, err := BookRepo.FindByID(id)
		if err != nil {
			http.Error(w, "Kitap bulunamadı", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(book)
		return
	}

	books, err := BookRepo.FindAll()
	if err != nil {
		http.Error(w, "Kitaplar alınamadı", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(books)
}
