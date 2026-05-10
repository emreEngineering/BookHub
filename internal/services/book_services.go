package services

import (
	"BookHub/internal/models"
	"BookHub/internal/repositories"
)

type BookServices interface {
	GetAllBooks() ([]models.Book, error)
	GetBookByID(id int) (*models.Book, error)
	CreateBook(book models.Book) (models.Book, error)
	UpdateBook(id int, book models.Book) (models.Book, error)
	DeleteBook(id int) error
}

type DefaultBookService struct {
	bookRepo repositories.BookRepository
}

func NewBookService(bookRepo repositories.BookRepository) *DefaultBookService {
	return &DefaultBookService{
		bookRepo: bookRepo,
	}
}
