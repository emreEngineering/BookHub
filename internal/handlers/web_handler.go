package handlers

import (
	"BookHub/internal/models"
	"BookHub/internal/responses"
	"BookHub/internal/services"
	"html/template"
	"net/http"
	"strconv"
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

	currentUser, isAuthenticated := h.currentUser(r)

	data := struct {
		Title           string
		Books           interface{}
		IsAuthenticated bool
		CurrentUser     *models.User
	}{
		Title:           "BookHub - Kitaplar",
		Books:           books,
		IsAuthenticated: isAuthenticated,
		CurrentUser:     currentUser,
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
		Title           string
		Error           string
		IsAuthenticated bool
		CurrentUser     *models.User
	}{
		Title:           "BookHub - Giriş Yap",
		Error:           errorMessage,
		IsAuthenticated: false,
		CurrentUser:     nil,
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
		Title           string
		Error           string
		Message         string
		IsAuthenticated bool
		CurrentUser     *models.User
	}{
		Title:           "BookHub - Kayıt Ol",
		Error:           errorMessage,
		Message:         "",
		IsAuthenticated: false,
		CurrentUser:     nil,
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Template çalıştırılamadı")
		return
	}
}

func (h *WebHandler) currentUser(r *http.Request) (*models.User, bool) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, false
	}

	userID, err := h.sessionService.GetUserID(cookie.Value)
	if err != nil {
		return nil, false
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		return nil, false
	}
	return user, true
}

func (h *WebHandler) LogoutPageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responses.Error(w, http.StatusMethodNotAllowed, "Bu endpoint sadece POST destekler")
		return
	}

	cookie, err := r.Cookie("session_id")
	if err == nil {
		h.sessionService.DeleteSession(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
	http.Redirect(w, r, "/web/login", http.StatusSeeOther)
}

func (h *WebHandler) BookCreatePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.renderBookFormPage(w, r, "")
		return
	}
	if r.Method == http.MethodPost {
		year, err := strconv.Atoi(r.FormValue("year"))
		if err != nil {
			h.renderBookFormPage(w, r, "Yıl geçerli bir sayı olmalıdır")
			return
		}

		book := models.Book{
			Title:  r.FormValue("title"),
			Author: r.FormValue("author"),
			Year:   year,
		}

		_, err = h.bookService.CreateBook(book)
		if err != nil {
			h.renderBookFormPage(w, r, err.Error())
			return
		}
		http.Redirect(w, r, "/web/books", http.StatusSeeOther)
		return
	}
	responses.Error(w, http.StatusMethodNotAllowed, "Bu endpoint sadece GET ve POST destekler")
}

func (h *WebHandler) renderBookFormPage(w http.ResponseWriter, r *http.Request, errorMessage string) {
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/book_form.html")
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Template yüklenemedi")
		return
	}

	currentUser, isAuthenticated := h.currentUser(r)

	data := struct {
		Title           string
		Error           string
		IsAuthenticated bool
		CurrentUser     *models.User
	}{
		Title:           "BookHub - Yeni Kitap Ekle",
		Error:           errorMessage,
		IsAuthenticated: isAuthenticated,
		CurrentUser:     currentUser,
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Template çalıştırılamadı")
		return
	}
}

func (h *WebHandler) BookDeletePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responses.Error(w, http.StatusMethodNotAllowed, "Bu endpoint sadece POST destekler")
		return
	}

	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		responses.Error(w, http.StatusBadRequest, "Kitap ID zorunludur")
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "Geçersiz kita ID")
		return
	}
	err = h.bookService.DeleteBook(id)
	if err != nil {
		responses.Error(w, http.StatusNotFound, "Kitap bulunamadı")
		return
	}
	http.Redirect(w, r, "/web/books", http.StatusSeeOther)
}

func (h *WebHandler) BookEditPageHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		responses.Error(w, http.StatusBadRequest, "Kitap ID zorunludur")
		return
	}
	id, err := strconv.Atoi(idParam)
	if err != nil {
		responses.Error(w, http.StatusBadRequest, "Geçersiz kitap ID")
		return
	}

	if r.Method == http.MethodGet {
		book, err := h.bookService.GetBookByID(id)
		if err != nil {
			responses.Error(w, http.StatusNotFound, "Kitap bulunamadı")
			return
		}
		h.renderBookEditPage(w, r, *book, "")
		return
	}

	if r.Method == http.MethodPost {
		year, err := strconv.Atoi(r.FormValue("year"))
		if err != nil {
			book := models.Book{
				ID:     id,
				Title:  r.FormValue("title"),
				Author: r.FormValue("author"),
				Year:   0,
			}
			h.renderBookEditPage(w, r, book, "Yıl geçerli bir sayı olamlıdır")
			return
		}

		book := models.Book{
			Title:  r.FormValue("title"),
			Author: r.FormValue("author"),
			Year:   year,
		}

		_, err = h.bookService.UpdateBook(id, book)
		if err != nil {
			book.ID = id
			h.renderBookEditPage(w, r, book, err.Error())
			return
		}
		http.Redirect(w, r, "/web/books", http.StatusSeeOther)
		return
	}
	responses.Error(w, http.StatusMethodNotAllowed, "Bu endpoint sadece GET ve POST destekler")
}

func (h *WebHandler) renderBookEditPage(w http.ResponseWriter, r *http.Request, book models.Book, errorMessage string) {
	tmpl, err := template.ParseFiles("templates/layout.html", "templates/book_edit.html")
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Template yüklenemedi")
		return
	}
	currentUser, isAuthenticated := h.currentUser(r)

	data := struct {
		Title           string
		Error           string
		Book            models.Book
		IsAuthenticated bool
		CurrentUser     *models.User
	}{
		Title:           "BookHub - Kitap Düzenle",
		Error:           errorMessage,
		Book:            book,
		IsAuthenticated: isAuthenticated,
		CurrentUser:     currentUser,
	}
	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		responses.Error(w, http.StatusInternalServerError, "Template çalıştırılamadı")
		return
	}
}
