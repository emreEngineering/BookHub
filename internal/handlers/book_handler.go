package handlers

import (
	"BookHub/internal/responses"
	"BookHub/internal/services"
	"encoding/json"
	"net/http"
	"strconv"

	"BookHub/internal/models"
)

type BookHandler struct {
	bookService services.BookServices
}

func NewBookHandler(bookService services.BookServices) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

func (h *BookHandler) BooksHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		var book models.Book

		err := json.NewDecoder(r.Body).Decode(&book)
		if err != nil {
			responses.Error(w, http.StatusBadRequest, "Geçersiz JSON")
			return
		}

		createdBook, err := h.bookService.CreateBook(book)
		if err != nil {
			responses.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		responses.Success(w, http.StatusCreated, "Kitap oluşturuldu", createdBook)
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
			if err.Error() == "kitap bulunamadı" {
				responses.Error(w, http.StatusNotFound, "Kitap bulunamadı")
				return
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		responses.Success(w, http.StatusOK, "Kitap güncellendi", updatedBook)
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

		responses.Success(w, http.StatusOK, "Kitap silindi", nil)
		return
	}

	if r.Method != http.MethodGet {
		responses.Error(w, http.StatusMethodNotAllowed, "Bu endpoint sadece GET, POST, PUT ve DELETE destekler")
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

		responses.Success(w, http.StatusOK, "Kitap getirildi", book)
		return
	}

	books, err := h.bookService.GetAllBooks()
	if err != nil {
		http.Error(w, "Kitaplar alınamadı", http.StatusInternalServerError)
		return
	}

	responses.Success(w, http.StatusOK, "Kitaplar listelendi", books)
}
