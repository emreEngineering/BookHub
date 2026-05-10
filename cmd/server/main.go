package main

import (
	"BookHub/internal/models"
	"BookHub/internal/repositories"
	"encoding/json"
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
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		var book models.Book

		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, "Gecersiz JSON", http.StatusBadRequest)
			return
		}

		if book.Title == "" {
			http.Error(w, "Kitap adi bos olamaz", http.StatusBadRequest)
			return
		}

		if book.Author == "" {
			http.Error(w, "Yazar adi bos olamaz", http.StatusBadRequest)
			return
		}

		createdBook, err := bookRepo.Create(book)
		if err != nil {
			http.Error(w, "Kitap olusturulamadi", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdBook)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Bu endpoint sadece GET ve POST destekler", http.StatusMethodNotAllowed)
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

		json.NewEncoder(w).Encode(book)
		return
	}

	books, err := bookRepo.FindAll()
	if err != nil {
		http.Error(w, "Kitaplar alinamadi", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(books)
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
