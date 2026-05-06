// Bu dosyanın çalıştırılabilir ana paket olduğunu belirtir.
package main

// Dış paketleri kullanmak için import bloğunu başlatır.
import (
	// Book modelini kullanmak için models paketini içe aktarır.
	"BookHub/internal/models"
	// Kitap repository yapısını kullanmak için repositories paketini içe aktarır.
	"BookHub/internal/repositories"
	// JSON verisini okumak ve yazmak için encoding/json paketini içe aktarır.
	"encoding/json"
	// Ekrana ve HTTP cevabına yazı yazmak için fmt paketini içe aktarır.
	"fmt"
	// HTTP sunucusu ve handler yazmak için net/http paketini içe aktarır.
	"net/http"
	// Yazıyı sayıya çevirmek için strconv paketini içe aktarır.
	"strconv"
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

// booksHandler, kitap listeleme, kitap bulma ve kitap ekleme işlemlerini yönetir.
func booksHandler(w http.ResponseWriter, r *http.Request) {
	// Cevabın JSON formatında olduğunu belirtir.
	w.Header().Set("Content-Type", "application/json")

	// Gelen istek POST ise yeni kitap ekleme işlemi yapılır.
	if r.Method == http.MethodPost {
		// İstekten okunacak kitabı tutmak için boş Book değişkeni oluşturur.
		var book models.Book

		// İstek gövdesindeki JSON verisini book değişkenine çevirir.
		err := json.NewDecoder(r.Body).Decode(&book)
		// JSON okunamazsa hata cevabı gönderir.
		if err != nil {
			// Kullanıcıya geçersiz JSON hatası döndürür.
			http.Error(w, "Geçersiz JSON", http.StatusBadRequest)
			// Hata sonrası fonksiyondan çıkar.
			return
		}

		// Kitap adı boş mu kontrol eder.
		if book.Title == "" {
			// Kitap adı boşsa kullanıcıya hata cevabı gönderir.
			http.Error(w, "Kitap adı boş olamaz", http.StatusBadRequest)
			// Hata sonrası fonksiyondan çıkar.
			return
		}

		// Yazar adı boş mu kontrol eder.
		if book.Author == "" {
			// Yazar adı boşsa kullanıcıya hata cevabı gönderir.
			http.Error(w, "Yazar adı boş olamaz", http.StatusBadRequest)
			// Hata sonrası fonksiyondan çıkar.
			return
		}

		// Doğrulanan kitabı repository üzerinden kaydeder.
		createdBook, err := bookRepo.Create(book)
		// Kitap kaydedilemezse hata kontrolü yapar.
		if err != nil {
			// Kullanıcıya sunucu hatası cevabı gönderir.
			http.Error(w, "Kitap oluşturulamadı", http.StatusInternalServerError)
			// Hata sonrası fonksiyondan çıkar.
			return
		}

		// Başarılı ekleme için HTTP 201 durum kodunu yazar.
		w.WriteHeader(http.StatusCreated)
		// Oluşturulan kitabı JSON olarak kullanıcıya gönderir.
		json.NewEncoder(w).Encode(createdBook)
		// POST işlemi tamamlandığı için fonksiyondan çıkar.
		return
	}

	if r.Method == http.MethodPut {
		idParam := r.URL.Query().Get("id")
		if idParam == "" {
			http.Error(w, "Kitap ID zorunludur.", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "Geçersiz kitap ID", http.StatusBadRequest)
			return
		}

		var book models.Book
		err = json.NewDecoder(r.Body).Decode(&book)
		if err != nil {
			http.Error(w, "Geçersiz JSON", http.StatusBadRequest)
			return
		}

		if book.Title == "" {
			http.Error(w, "Kitap adı boş olamaz", http.StatusBadRequest)
			return
		}

		if book.Author == "" {
			http.Error(w, "Yazar adı boş olaamaz", http.StatusBadRequest)
			return
		}

		updatedBook, err := bookRepo.Update(id, book)

		if err != nil {
			http.Error(w, "Kitap bulunamadı", http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(updatedBook)
	}

	if r.Method == http.MethodDelete {
		idParam := r.URL.Query().Get("id")
		if idParam == "" {
			http.Error(w, "Kitap ID zorunludur", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "Geçersiz kitap ID", http.StatusBadRequest)
			return
		}

		err = bookRepo.Delete(id)
		if err != nil {
			http.Error(w, "Kitap bulunamadı", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}

	// İstek GET değilse bu endpoint için izin verilmez.
	if r.Method != http.MethodGet {
		// Kullanıcıya desteklenen metotları belirten hata cevabı gönderir.
		http.Error(w, "Bu endpoint sadece GET, POST, PUT ve DELETE destekler", http.StatusMethodNotAllowed)
		// Hata sonrası fonksiyondan çıkar.
		return
	}

	// URL query içindeki id değerini okur.
	idParam := r.URL.Query().Get("id")

	// id parametresi doluysa tek kitap arama işlemi yapılır.
	if idParam != "" {
		// id değerini yazıdan sayıya çevirir.
		id, err := strconv.Atoi(idParam)
		// id sayıya çevrilemezse hata kontrolü yapar.
		if err != nil {
			// Kullanıcıya geçersiz ID hatası gönderir.
			http.Error(w, "Geçersiz kitap ID", http.StatusBadRequest)
			// Hata sonrası fonksiyondan çıkar.
			return
		}

		// Repository içinde ID ile kitabı arar.
		book, err := bookRepo.FindByID(id)
		// Kitap bulunamazsa hata kontrolü yapar.
		if err != nil {
			// Kullanıcıya kitap bulunamadı cevabı gönderir.
			http.Error(w, "Kitap bulunamadı", http.StatusNotFound)
			// Hata sonrası fonksiyondan çıkar.
			return
		}

		// Bulunan kitabı JSON olarak kullanıcıya gönderir.
		json.NewEncoder(w).Encode(book)
		// Tek kitap arama işlemi tamamlandığı için fonksiyondan çıkar.
		return
	}

	// id parametresi yoksa bütün kitapları repository'den alır.
	books, err := bookRepo.FindAll()
	// Kitaplar alınamazsa hata kontrolü yapar.
	if err != nil {
		// Kullanıcıya sunucu hatası cevabı gönderir.
		http.Error(w, "Kitaplar alınamadı", http.StatusInternalServerError)
		// Hata sonrası fonksiyondan çıkar.
		return
	}

	// Bütün kitapları JSON olarak kullanıcıya gönderir.
	json.NewEncoder(w).Encode(books)
}

// main, uygulamanın başlangıç fonksiyonudur.
func main() {
	// Bellekte çalışan kitap deposunu oluşturur ve global değişkene atar.
	bookRepo = repositories.NewMemoryBookRepository()

	// Ana sayfa adresini homeHandler fonksiyonuna bağlar.
	http.HandleFunc("/", homeHandler)
	// Sağlık kontrol adresini healthHandler fonksiyonuna bağlar.
	http.HandleFunc("/health", healthHandler)
	// Hakkında sayfası adresini aboutHandler fonksiyonuna bağlar.
	http.HandleFunc("/about", aboutHandler)
	// Kitap işlemleri adresini booksHandler fonksiyonuna bağlar.
	http.HandleFunc("/books", booksHandler)

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
