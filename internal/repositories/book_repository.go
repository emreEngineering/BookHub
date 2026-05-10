package repositories

import (
	"BookHub/internal/models"
	"errors"
)

type BookRepository interface {
	FindAll() ([]models.Book, error)
	FindByID(id int) (*models.Book, error)
	Create(book models.Book) (models.Book, error)
	Update(id int, book models.Book) (models.Book, error)
	Delete(id int) error
}

type MemoryBookRepository struct {
	books  []models.Book
	nextID int
}

func NewMemoryBookRepository() *MemoryBookRepository {
	return &MemoryBookRepository{
		nextID: 4,
		books: []models.Book{
			{ID: 1, Title: "Suc ve Ceza", Author: "Dostoyevski", Year: 1866},
			{ID: 2, Title: "1984", Author: "George Orwell", Year: 1949},
			{ID: 3, Title: "Kurk Mantolu Madonna", Author: "Sabahattin Ali", Year: 1943},
		},
	}
}

func (r *MemoryBookRepository) FindAll() ([]models.Book, error) {
	return r.books, nil
}

func (r *MemoryBookRepository) FindByID(id int) (*models.Book, error) {
	for _, book := range r.books {
		if book.ID == id {
			found := book
			return &found, nil
		}
	}

	return nil, errors.New("book not found")
}

func (r *MemoryBookRepository) Create(book models.Book) (models.Book, error) {
	book.ID = r.nextID
	r.nextID++
	r.books = append(r.books, book)

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

	return models.Book{}, errors.New("book not found")
}

func (r *MemoryBookRepository) Delete(id int) error {
	for i, book := range r.books {
		if book.ID == id {
			r.books = append(r.books[:i], r.books[i+1:]...)
			return nil
		}
	}

	return errors.New("book not found")
}
