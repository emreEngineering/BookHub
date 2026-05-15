package handlers

import (
	"BookHub/internal/responses"
	"BookHub/internal/services"
	"html/template"
	"net/http"
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

func (h *WebHandler) LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responses.Error(w, http.StatusMethodNotAllowed, "Bu endpoint sadece GET destekler")
		return
	}

	tmpl, err := template.ParseFiles("templates/layout.html", "templates/login.html")
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Template yüklenemedi")
		return
	}

	data := struct {
		Title string
	}{
		Title: "BookHub - Giriş Yap",
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Template çalıştırılamadı")
		return
	}
}

func (h *WebHandler) RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responses.Error(w, http.StatusMethodNotAllowed, "Bu endpoint sadece GET destekler")
		return
	}

	tmpl, err := template.ParseFiles("templates/layout.html", "templates/register.html")
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Template yüklenemedi")
		return
	}

	data := struct {
		Title string
	}{
		Title: "BookHub - Kayıt Ol",
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Template çalıştırılamadı")
		return
	}
}
