package services

import (
	"context"
	"errors"

	"BookHub/internal/activity"
	"BookHub/internal/models"
	"BookHub/internal/repositories"
	"BookHub/internal/requestcontext"
)

type BookServices interface {
	GetAllBooks(ctx context.Context) ([]models.Book, error)
	GetBookByID(ctx context.Context, id int) (*models.Book, error)
	CreateBook(ctx context.Context, book models.Book) (models.Book, error)
	UpdateBook(ctx context.Context, id int, book models.Book) (models.Book, error)
	DeleteBook(ctx context.Context, id int) error
}

type DefaultBookService struct {
	bookRepo       repositories.BookRepository
	activityLogger activity.ActivityLogger
}

func NewBookService(bookRepo repositories.BookRepository, activityLoggers ...activity.ActivityLogger) *DefaultBookService {
	service := &DefaultBookService{
		bookRepo: bookRepo,
	}

	if len(activityLoggers) > 0 {
		service.activityLogger = activityLoggers[0]
	}

	return service
}

func (s *DefaultBookService) GetAllBooks(ctx context.Context) ([]models.Book, error) {
	return s.bookRepo.FindAll(ctx)
}
func (s *DefaultBookService) GetBookByID(ctx context.Context, id int) (*models.Book, error) {
	return s.bookRepo.FindByID(ctx, id)
}
func (s *DefaultBookService) CreateBook(ctx context.Context, book models.Book) (models.Book, error) {
	err := validateBook(book)
	if err != nil {
		return models.Book{}, err
	}

	createdBook, err := s.bookRepo.Create(ctx, book)
	if err != nil {
		return models.Book{}, err
	}

	s.logActivity(ctx, "book_created", "Kitap oluşturuldu", map[string]interface{}{
		"book_id": createdBook.ID,
		"title":   createdBook.Title,
		"author":  createdBook.Author,
		"year":    createdBook.Year,
	})

	return createdBook, nil
}
func (s *DefaultBookService) UpdateBook(ctx context.Context, id int, book models.Book) (models.Book, error) {
	err := validateBook(book)

	if err != nil {
		return models.Book{}, err
	}
	updatedBook, err := s.bookRepo.Update(ctx, id, book)
	if err != nil {
		return models.Book{}, err
	}

	s.logActivity(ctx, "book_updated", "Kitap güncellendi", map[string]interface{}{
		"book_id": updatedBook.ID,
		"title":   updatedBook.Title,
		"author":  updatedBook.Author,
		"year":    updatedBook.Year,
	})

	return updatedBook, nil
}

func (s *DefaultBookService) DeleteBook(ctx context.Context, id int) error {
	err := s.bookRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	s.logActivity(ctx, "book_deleted", "Kitap silindi", map[string]interface{}{
		"book_id": id,
	})

	return nil
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

func (s *DefaultBookService) logActivity(ctx context.Context, eventType string, message string, metadata map[string]interface{}) {
	if s.activityLogger == nil {
		return
	}

	var userIDPtr *int
	if userID, ok := requestcontext.UserID(ctx); ok {
		userIDPtr = &userID
	}

	_ = s.activityLogger.Log(ctx, eventType, message, userIDPtr, metadata)
}
