package main

import (
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
	bookHandler := handlers.NewBookHandler(bookRepo)

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/books", bookHandler.BooksHandler)

	fmt.Println("Server çalışıyor: http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server başlatılamadı:", err)
	}
}
