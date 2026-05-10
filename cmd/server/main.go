package main

import (
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
	authHandler := handlers.NewAuthHandler(userService)
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/books", bookHandler.BooksHandler)
	http.HandleFunc("/register", authHandler.RegisterHandler)

	fmt.Println("Server çalışıyor: http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server başlatılamadı:", err)
	}
}
