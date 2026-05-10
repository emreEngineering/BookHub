// Bu dosyanın çalıştırılabilir ana paket olduğunu belirtir.
package main

// Dış paketleri kullanmak için import bloğunu başlatır.
import (
	"BookHub/internal/handlers"
	// Kitap repository yapısını kullanmak için repositories paketini içe aktarır.
	"BookHub/internal/repositories"
	// Ekrana ve HTTP cevabına yazı yazmak için fmt paketini içe aktarır.
	"fmt"
	// HTTP sunucusu ve handler yazmak için net/http paketini içe aktarır.
	"net/http"
)

// bookRepo, bütün handler fonksiyonlarının kullanacağı ortak kitap deposudur.
var bookRepo repositories.BookRepository

// homeHandler, ana sayfa isteğine cevap verir.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Kullanıcıya kısa bir çalışma mesajı gönderir.
	fmt.Fprintln(w, "BookHub çalışıyor.")
}

// healthHandler, uygulamanın ayakta olup olmadığını bildirir.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	// Sağlık kontrolü için OK cevabı gönderir.
	fmt.Fprintln(w, "OK")
}

// aboutHandler, uygulama hakkında kısa bilgi verir.
func aboutHandler(w http.ResponseWriter, r *http.Request) {
	// Kullanıcıya uygulamanın açıklamasını gönderir.
	fmt.Fprintln(w, "BookHub bir kitap yönetim sistemidir.")
}

// main, uygulamanın başlangıç fonksiyonudur.
func main() {
	// Bellekte çalışan kitap deposunu oluşturur ve global değişkene atar.
	bookRepo = repositories.NewMemoryBookRepository()
	handlers.BookRepo = bookRepo
	bookHandler := handlers.NewBookHandler()

	// Ana sayfa adresini homeHandler fonksiyonuna bağlar.
	http.HandleFunc("/", homeHandler)
	// Sağlık kontrol adresini healthHandler fonksiyonuna bağlar.
	http.HandleFunc("/health", healthHandler)
	// Hakkında sayfası adresini aboutHandler fonksiyonuna bağlar.
	http.HandleFunc("/about", aboutHandler)
	// Kitap işlemleri adresini booksHandler fonksiyonuna bağlar.
	http.HandleFunc("/books", bookHandler.BooksHandler)

	// Sunucunun hangi adreste çalıştığını terminale yazar.
	fmt.Println("Server çalışıyor: http://localhost:8080")

	// HTTP sunucusunu 8080 portunda başlatır.
	err := http.ListenAndServe(":8080", nil)

	// Sunucu başlatılırken hata oluşup oluşmadığını kontrol eder.
	if err != nil {
		// Hata varsa terminale hata mesajını yazar.
		fmt.Println("Server başlatılamadı:", err)
	}
}
