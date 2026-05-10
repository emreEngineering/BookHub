package services

import "BookHub/internal/models"

type BookServices interface {
	GetAllBooks() ([]models.Book, error)
	GetBookByID(id int) (*models.Book, error)
	CreateBook(book models.Book) (models.Book, error)
	UpdateBook(id int, book models.Book) (models.Book, error)
	DeleteBook(id int) error
}
