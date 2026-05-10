package services

import (
	"BookHub/internal/models"
	"BookHub/internal/repositories"
	"errors"
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

func (s *DefaultBookService) GetAllBooks() ([]models.Book, error) {
	return s.bookRepo.FindAll()
}

func (s *DefaultBookService) GetBookByID(id int) (*models.Book, error) {
	return s.bookRepo.FindByID(id)
}

func (s *DefaultBookService) CreateBook(book models.Book) (models.Book, error) {
	if err := validateBook(book); err != nil {
		return models.Book{}, err
	}

	return s.bookRepo.Create(book)
}

func (s *DefaultBookService) UpdateBook(id int, book models.Book) (models.Book, error) {
	if err := validateBook(book); err != nil {
		return models.Book{}, err
	}

	return s.bookRepo.Update(id, book)
}

func (s *DefaultBookService) DeleteBook(id int) error {
	return s.bookRepo.Delete(id)
}

func validateBook(book models.Book) error {
	if book.Title == "" {
		return errors.New("Kitap adı boş olamaz")
	}

	if book.Author == "" {
		return errors.New("Yazar adı boş olamaz")
	}

	return nil
}
