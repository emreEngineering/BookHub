// Bu dosyanın repositories paketine ait olduğunu belirtir.
package repositories

// Dış paketleri kullanmak için import bloğunu başlatır.
import (
	// Book modelini kullanmak için models paketini içe aktarır.
	"BookHub/internal/models"
	// Hata mesajı oluşturmak için errors paketini içe aktarır.
	"errors"
)

// BookRepository, kitap verileri için gerekli işlemleri tanımlar.
type BookRepository interface {
	// FindAll, bütün kitapları döndürür.
	FindAll() ([]models.Book, error)
	// FindByID, verilen ID'ye sahip kitabı döndürür.
	FindByID(id int) (*models.Book, error)
	// Create, yeni kitabı kaydeder ve döndürür.
	Create(book models.Book) (models.Book, error)
	//Şu ID’ye sahip kitabı, gönderilen yeni kitap bilgileriyle güncelle.
	Update(id int, book models.Book) (models.Book, error)
	Delete(id int) error
}

// MemoryBookRepository, kitapları bellekte tutan repository yapısıdır.
type MemoryBookRepository struct {
	// books, bellekte saklanan kitap listesidir.
	books []models.Book
	// nextID, yeni kitap için kullanılacak sıradaki ID değeridir.
	nextID int
}

// NewMemoryBookRepository, başlangıç kitaplarıyla yeni repository oluşturur.
func NewMemoryBookRepository() *MemoryBookRepository {
	// Bellekte çalışan yeni repository adresini döndürür.
	return &MemoryBookRepository{
		// Yeni eklenecek ilk kitabın ID değeri 5 olur.
		nextID: 5,
		// Başlangıçta sistemde bulunacak kitap listesini tanımlar.
		books: []models.Book{
			// Birinci örnek kitabı tanımlar.
			{
				// Kitabın ID değerini belirler.
				ID: 1,
				// Kitabın adını belirler.
				Title: "Suç ve Ceza",
				// Kitabın yazarını belirler.
				Author: "Dostoyevski",
				// Kitabın yayın yılını belirler.
				Year: 1866,
			},
			// İkinci örnek kitabı tanımlar.
			{
				// Kitabın ID değerini belirler.
				ID: 2,
				// Kitabın adını belirler.
				Title: "1984",
				// Kitabın yazarını belirler.
				Author: "George Orwell",
				// Kitabın yayın yılını belirler.
				Year: 1949,
			},
			// Üçüncü örnek kitabı tanımlar.
			{
				// Kitabın ID değerini belirler.
				ID: 3,
				// Kitabın adını belirler.
				Title: "Kürk Mantolu Madonna",
				// Kitabın yazarını belirler.
				Author: "Sabahattin Ali",
				// Kitabın yayın yılını belirler.
				Year: 1943,
			},
			// Dördüncü örnek kitabı tanımlar.
			{
				// Kitabın ID değerini belirler.
				ID: 4,
				// Kitabın adını belirler.
				Title: "Tutunamayanlar",
				// Kitabın yazarını belirler.
				Author: "Oğuz Atay",
				// Kitabın yayın yılını belirler.
				Year: 1972,
			},
		},
	}
}

// FindAll, bellekteki tüm kitap listesini döndürür.
func (r *MemoryBookRepository) FindAll() ([]models.Book, error) {
	// Kitap listesini ve hata olmadığını belirtmek için nil döndürür.
	return r.books, nil
}

// FindByID, verilen ID ile eşleşen kitabı arar.
func (r *MemoryBookRepository) FindByID(id int) (*models.Book, error) {
	// Kitap listesindeki her kitap için döngü başlatır.
	for _, book := range r.books {
		// Kitabın ID değeri aranan ID ile aynı mı kontrol eder.
		if book.ID == id {
			// Eşleşen kitabın adresini ve hata olmadığını döndürür.
			return &book, nil
		}
	}
	// Kitap bulunamazsa nil ve hata mesajı döndürür.
	return nil, errors.New("Kitap bulunamadı.")
}

// Create, yeni kitabı bellekteki listeye ekler.
func (r *MemoryBookRepository) Create(book models.Book) (models.Book, error) {
	// Yeni kitaba sıradaki ID değerini verir.
	book.ID = r.nextID
	// Bir sonraki kitap için ID değerini artırır.
	r.nextID++

	// Yeni kitabı kitap listesine ekler.
	r.books = append(r.books, book)

	// Eklenen kitabı ve hata olmadığını belirtmek için nil döndürür.
	return book, nil
}

func (r *MemoryBookRepository) Update(id int, book models.Book) (models.Book, error) {
	for i, existingBook := range r.books {
		if existingBook.ID == id {
			book.ID = id
			r.books[i] = book
			return book, nil
		}
	}
	return models.Book{}, errors.New("Kitap bulunamadı")
}

func (r *MemoryBookRepository) Delete(id int) error {
	for i, book := range r.books {
		if book.ID == id {
			r.books = append(r.books[:i], r.books[i+1:]...)
			return nil
		}
	}
	return errors.New("Kitap Bulunamadı")
}
