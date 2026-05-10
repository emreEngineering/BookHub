package repositories

import (
	"BookHub/internal/models"
	"errors"
)

type BookRepository interface {
	FindAll() ([]models.Book, error)
	FindByID(id int) (*models.Book, error)
}

type MemoryBookRepository struct {
	books []models.Book
}

func NewMemoryBookRepository() *MemoryBookRepository {
	return &MemoryBookRepository{
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
