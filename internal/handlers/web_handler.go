package handlers

import (
	"BookHub/internal/responses"
	"BookHub/internal/services"
	"html/template"
	"net/http"
)

type WebHandler struct {
	bookService    services.BookServices
	userService    services.UserService
	sessionService services.SessionService
}

func NewWebHandler(bookService services.BookServices, userService services.UserService, sessionService services.SessionService) *WebHandler {
	return &WebHandler{
		bookService:    bookService,
		userService:    userService,
		sessionService: sessionService,
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
	if r.Method == http.MethodGet {
		h.renderLoginPage(w, "")
		return
	}

	if r.Method == http.MethodPost {
		request := services.LoginRequest{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		user, err := h.userService.Login(request)
		if err != nil {
			h.renderLoginPage(w, err.Error())
			return
		}

		sessionID, err := h.sessionService.CreateSession(user.ID)
		if err != nil {
			h.renderLoginPage(w, "Session oluşturulamadı")
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   3600,
		})

		http.Redirect(w, r, "/web/books", http.StatusSeeOther)
		return
	}

	responses.Error(w, http.StatusMethodNotAllowed, "Bu endpoint sadece GET ve POST destekler")
}

func (h *WebHandler) renderLoginPage(w http.ResponseWriter, errorMessage string) {
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/login.html")
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Template yüklenemedi")
		return
	}

	data := struct {
		Title string
		Error string
	}{
		Title: "BookHub - Giriş Yap",
		Error: errorMessage,
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Template çalıştırılamadı")
		return
	}
}

func (h *WebHandler) RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.renderRegisterPage(w, "")
		return
	}

	if r.Method == http.MethodPost {
		request := services.RegisterRequest{
			Name:     r.FormValue("name"),
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
		}

		_, err := h.userService.Register(request)
		if err != nil {
			h.renderRegisterPage(w, err.Error())
			return
		}

		http.Redirect(w, r, "/web/login", http.StatusSeeOther)
		return
	}

	responses.Error(w, http.StatusMethodNotAllowed, "Bu endpoint sadece GET ve POST destekler")
}
func (h *WebHandler) renderRegisterPage(w http.ResponseWriter, errorMessage string) {
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/register.html")
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Template yüklenemedi")
		return
	}

	data := struct {
		Title   string
		Error   string
		Message string
	}{
		Title:   "BookHub - Kayıt Ol",
		Error:   errorMessage,
		Message: "",
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Template çalıştırılamadı")
		return
	}
}
