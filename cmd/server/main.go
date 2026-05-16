package main

import (
	"BookHub/internal/middleware"
	"BookHub/internal/services"
	"fmt"
	"net/http"

	"BookHub/internal/handlers"
	"BookHub/internal/repositories"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "BookHub çalışıyor")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "BookHub bir kitap yönetim sistemidir")
}

func main() {
	bookRepo := repositories.NewMemoryBookRepository()
	bookService := services.NewBookService(bookRepo)
	bookHandler := handlers.NewBookHandler(bookService)

	userRepo := repositories.NewMemoryUserRepository()
	userService := services.NewUserService(userRepo)
	sessionService := services.NewSessionService()
	authHandler := handlers.NewAuthHandler(userService, sessionService)
	authMiddleware := middleware.NewAuthMiddleware(sessionService)

	webHandler := handlers.NewWebHandler(bookService, userService, sessionService)

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/about", aboutHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/books", bookHandler.BooksHandler)
	http.HandleFunc("/web/books", webHandler.BooksPageHandler)
	http.HandleFunc("/web/login", webHandler.LoginPageHandler)
	http.HandleFunc("/web/register", webHandler.RegisterPageHandler)
	http.HandleFunc("/register", authHandler.RegisterHandler)
	http.HandleFunc("/login", authHandler.LoginHandler)
	http.HandleFunc("/me", authMiddleware.RequireAuth(authHandler.MeHandler))
	http.HandleFunc("/logout", authMiddleware.RequireAuth(authHandler.LogoutHandler))
	fmt.Println("Server çalışıyor: http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server başlatılamadı:", err)
	}
}
