package repositories

import (
	"context"

	"BookHub/internal/models"
)

type BookRepository interface {
	FindAll(ctx context.Context) ([]models.Book, error)
	FindByID(ctx context.Context, id int) (*models.Book, error)
	Create(ctx context.Context, book models.Book) (models.Book, error)
	Update(ctx context.Context, id int, book models.Book) (models.Book, error)
	Delete(ctx context.Context, id int) error
}
