package main

import (
	"BookHub/internal/repositories"
	"fmt"
	"net/http"
	"strconv"
)

var bookRepo repositories.BookRepository

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "BookHub calisiyor")
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "BookHub bir kitap yonetim sistemidir")
}

func booksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Bu endpoint sadece GET destekler", http.StatusMethodNotAllowed)
		return
	}

	idParam := r.URL.Query().Get("id")
	if idParam != "" {
		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "Gecersiz kitap ID", http.StatusBadRequest)
			return
		}

		book, err := bookRepo.FindByID(id)
		if err != nil {
			http.Error(w, "Kitap bulunamadi", http.StatusNotFound)
			return
		}

		fmt.Fprintf(w, "%d - %s / %s (%d)\n", book.ID, book.Title, book.Author, book.Year)
		return
	}

	books, err := bookRepo.FindAll()
	if err != nil {
		http.Error(w, "Kitaplar alinamadi", http.StatusInternalServerError)
		return
	}

	for _, book := range books {
		fmt.Fprintf(w, "%d - %s / %s (%d)\n", book.ID, book.Title, book.Author, book.Year)
	}
}

func main() {
	bookRepo = repositories.NewMemoryBookRepository()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/about", aboutHandler)
	http.HandleFunc("/books", booksHandler)

	fmt.Println("Server calisiyor: http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server baslatilamadi:", err)
	}
}
