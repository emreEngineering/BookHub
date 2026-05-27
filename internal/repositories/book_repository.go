package repositories

import "BookHub/internal/models"

type BookRepository interface {
	FindAll() ([]models.Book, error)
	FindByID(id int) (*models.Book, error)
	Create(book models.Book) (models.Book, error)
	Update(id int, book models.Book) (models.Book, error)
	Delete(id int) error
}
