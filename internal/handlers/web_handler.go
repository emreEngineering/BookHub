package handlers

import (
	"html/template"
	"net/http"

	"BookHub/internal/responses"
	"BookHub/internal/services"
)

type WebHandler struct {
	bookService services.BookServices
}

func NewWebHandler(bookService services.BookServices) *WebHandler {
	return &WebHandler{
		bookService: bookService,
	}
}

func (h *WebHandler) BooksPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responses.Error(w, http.StatusMethodNotAllowed, "Bu endpoint sadece GET destekler")
		return
	}

	books, err := h.bookService.GetAllBooks()
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Kitaplar alınamadı")
		return
	}

	tmpl, err := template.ParseFiles("templates/layout.html", "templates/books.html")
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Template yüklenemedi")
		return
	}

	data := struct {
		Title string
		Books interface{}
	}{
		Title: "BookHub - Kitaplar",
		Books: books,
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Template çalıştırılamadı")
		return
	}
}
