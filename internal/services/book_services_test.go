package services

import (
	"context"
	"errors"
	"testing"

	"BookHub/internal/activity"
	"BookHub/internal/models"
	"BookHub/internal/requestcontext"
)

type fakeBookRepository struct {
	books  []models.Book
	nextID int
}

type fakeBookActivityLogger struct {
	events []activity.ActivityEvent
}

func newFakeBookRepository() *fakeBookRepository {
	return &fakeBookRepository{
		nextID: 3,
		books: []models.Book{
			{ID: 1, Title: "Book One", Author: "Author One", Year: 2001},
			{ID: 2, Title: "Book Two", Author: "Author Two", Year: 2002},
		},
	}
}

func (r *fakeBookRepository) FindAll(ctx context.Context) ([]models.Book, error) {
	return r.books, nil
}

func (r *fakeBookRepository) FindByID(ctx context.Context, id int) (*models.Book, error) {
	for _, book := range r.books {
		if book.ID == id {
			return &book, nil
		}
	}

	return nil, errors.New("kitap bulunamadı")
}

func (r *fakeBookRepository) Create(ctx context.Context, book models.Book) (models.Book, error) {
	book.ID = r.nextID
	r.nextID++
	r.books = append(r.books, book)
	return book, nil
}

func (r *fakeBookRepository) Update(ctx context.Context, id int, book models.Book) (models.Book, error) {
	for i := range r.books {
		if r.books[i].ID == id {
			book.ID = id
			r.books[i] = book
			return book, nil
		}
	}

	return models.Book{}, errors.New("kitap bulunamadı")
}

func (r *fakeBookRepository) Delete(ctx context.Context, id int) error {
	for i := range r.books {
		if r.books[i].ID == id {
			r.books = append(r.books[:i], r.books[i+1:]...)
			return nil
		}
	}

	return errors.New("kitap bulunamadı")
}

func (l *fakeBookActivityLogger) Log(ctx context.Context, eventType string, message string, userID *int, metadata map[string]interface{}) error {
	l.events = append(l.events, activity.ActivityEvent{
		Type:     eventType,
		Message:  message,
		UserID:   userID,
		Metadata: metadata,
	})
	return nil
}

func (l *fakeBookActivityLogger) FindLatest(ctx context.Context, limit int64) ([]activity.ActivityLog, error) {
	return nil, nil
}

func TestBookService_GetAllBooks(t *testing.T) {
	service := NewBookService(newFakeBookRepository())

	books, err := service.GetAllBooks(context.Background())
	if err != nil {
		t.Fatalf("GetAllBooks returned error: %v", err)
	}

	if len(books) != 2 {
		t.Fatalf("expected 2 books, got %d", len(books))
	}
}

func TestBookService_GetBookByID_Found(t *testing.T) {
	service := NewBookService(newFakeBookRepository())

	book, err := service.GetBookByID(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetBookByID returned error: %v", err)
	}

	if book == nil || book.ID != 1 {
		t.Fatalf("expected book ID 1, got %#v", book)
	}
}

func TestBookService_GetBookByID_NotFound(t *testing.T) {
	service := NewBookService(newFakeBookRepository())

	book, err := service.GetBookByID(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error for missing book")
	}

	if book != nil {
		t.Fatalf("expected nil book, got %#v", book)
	}
}

func TestBookService_CreateBook_Success(t *testing.T) {
	service := NewBookService(newFakeBookRepository())

	book, err := service.CreateBook(context.Background(), models.Book{
		Title:  "New Book",
		Author: "New Author",
		Year:   2026,
	})
	if err != nil {
		t.Fatalf("CreateBook returned error: %v", err)
	}

	if book.ID == 0 {
		t.Fatal("expected created book ID to be set")
	}
}

func TestBookService_CreateBook_LogsActivityWithActor(t *testing.T) {
	logger := &fakeBookActivityLogger{}
	service := NewBookService(newFakeBookRepository(), logger)
	ctx := requestcontext.WithUserID(context.Background(), 42)

	book, err := service.CreateBook(ctx, models.Book{
		Title:  "New Book",
		Author: "New Author",
		Year:   2026,
	})
	if err != nil {
		t.Fatalf("CreateBook returned error: %v", err)
	}

	if len(logger.events) != 1 {
		t.Fatalf("expected 1 activity event, got %d", len(logger.events))
	}

	event := logger.events[0]
	if event.Type != "book_created" {
		t.Fatalf("expected book_created event, got %q", event.Type)
	}
	if event.UserID == nil || *event.UserID != 42 {
		t.Fatalf("expected userID 42, got %#v", event.UserID)
	}
	if event.Metadata["book_id"] != book.ID {
		t.Fatalf("expected metadata book_id %d, got %#v", book.ID, event.Metadata["book_id"])
	}
}

func TestBookService_CreateBook_ValidationError(t *testing.T) {
	service := NewBookService(newFakeBookRepository())

	_, err := service.CreateBook(context.Background(), models.Book{
		Title:  "",
		Author: "Author",
		Year:   2026,
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestBookService_UpdateBook_Success(t *testing.T) {
	service := NewBookService(newFakeBookRepository())

	book, err := service.UpdateBook(context.Background(), 1, models.Book{
		Title:  "Updated Book",
		Author: "Updated Author",
		Year:   2027,
	})
	if err != nil {
		t.Fatalf("UpdateBook returned error: %v", err)
	}

	if book.ID != 1 || book.Title != "Updated Book" {
		t.Fatalf("unexpected updated book: %#v", book)
	}
}

func TestBookService_UpdateBook_ValidationError(t *testing.T) {
	service := NewBookService(newFakeBookRepository())

	_, err := service.UpdateBook(context.Background(), 1, models.Book{
		Title:  "Updated Book",
		Author: "",
		Year:   2027,
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestBookService_DeleteBook_Success(t *testing.T) {
	repo := newFakeBookRepository()
	service := NewBookService(repo)

	err := service.DeleteBook(context.Background(), 1)
	if err != nil {
		t.Fatalf("DeleteBook returned error: %v", err)
	}

	if len(repo.books) != 1 {
		t.Fatalf("expected 1 book after delete, got %d", len(repo.books))
	}
}

func TestBookService_DeleteBook_NotFound(t *testing.T) {
	service := NewBookService(newFakeBookRepository())

	err := service.DeleteBook(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error for missing book")
	}
}
