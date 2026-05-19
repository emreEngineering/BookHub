package services

import (
	"errors"
	"testing"

	"BookHub/internal/models"
)

type fakeBookRepository struct {
	books  []models.Book
	nextID int
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

func (r *fakeBookRepository) FindAll() ([]models.Book, error) {
	return r.books, nil
}

func (r *fakeBookRepository) FindByID(id int) (*models.Book, error) {
	for _, book := range r.books {
		if book.ID == id {
			return &book, nil
		}
	}

	return nil, errors.New("kitap bulunamadı")
}

func (r *fakeBookRepository) Create(book models.Book) (models.Book, error) {
	book.ID = r.nextID
	r.nextID++
	r.books = append(r.books, book)
	return book, nil
}

func (r *fakeBookRepository) Update(id int, book models.Book) (models.Book, error) {
	for i := range r.books {
		if r.books[i].ID == id {
			book.ID = id
			r.books[i] = book
			return book, nil
		}
	}

	return models.Book{}, errors.New("kitap bulunamadı")
}

func (r *fakeBookRepository) Delete(id int) error {
	for i := range r.books {
		if r.books[i].ID == id {
			r.books = append(r.books[:i], r.books[i+1:]...)
			return nil
		}
	}

	return errors.New("kitap bulunamadı")
}

func TestBookService_GetAllBooks(t *testing.T) {
	service := NewBookService(newFakeBookRepository())

	books, err := service.GetAllBooks()
	if err != nil {
		t.Fatalf("GetAllBooks returned error: %v", err)
	}

	if len(books) != 2 {
		t.Fatalf("expected 2 books, got %d", len(books))
	}
}

func TestBookService_GetBookByID_Found(t *testing.T) {
	service := NewBookService(newFakeBookRepository())

	book, err := service.GetBookByID(1)
	if err != nil {
		t.Fatalf("GetBookByID returned error: %v", err)
	}

	if book == nil || book.ID != 1 {
		t.Fatalf("expected book ID 1, got %#v", book)
	}
}

func TestBookService_GetBookByID_NotFound(t *testing.T) {
	service := NewBookService(newFakeBookRepository())

	book, err := service.GetBookByID(999)
	if err == nil {
		t.Fatal("expected error for missing book")
	}

	if book != nil {
		t.Fatalf("expected nil book, got %#v", book)
	}
}

func TestBookService_CreateBook_Success(t *testing.T) {
	service := NewBookService(newFakeBookRepository())

	book, err := service.CreateBook(models.Book{
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

func TestBookService_CreateBook_ValidationError(t *testing.T) {
	service := NewBookService(newFakeBookRepository())

	_, err := service.CreateBook(models.Book{
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

	book, err := service.UpdateBook(1, models.Book{
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

	_, err := service.UpdateBook(1, models.Book{
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

	err := service.DeleteBook(1)
	if err != nil {
		t.Fatalf("DeleteBook returned error: %v", err)
	}

	if len(repo.books) != 1 {
		t.Fatalf("expected 1 book after delete, got %d", len(repo.books))
	}
}

func TestBookService_DeleteBook_NotFound(t *testing.T) {
	service := NewBookService(newFakeBookRepository())

	err := service.DeleteBook(999)
	if err == nil {
		t.Fatal("expected error for missing book")
	}
}
