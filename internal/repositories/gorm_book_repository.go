package repositories

import (
	gormmodels "BookHub/internal/gormdb/models"
	"BookHub/internal/models"
	"errors"

	"gorm.io/gorm"
)

type GormBookRepository struct {
	db *gorm.DB
}

func NewGormBookRepository(db *gorm.DB) *GormBookRepository {
	return &GormBookRepository{
		db: db,
	}
}

func (r *GormBookRepository) FindAll() ([]models.Book, error) {
	var gormBooks []gormmodels.GormBook

	err := r.db.Order("id").Find(&gormBooks).Error
	if err != nil {
		return nil, err
	}

	books := make([]models.Book, 0, len(gormBooks))
	for _, gormBook := range gormBooks {
		books = append(books, toBook(gormBook))
	}

	return books, nil
}

func (r *GormBookRepository) FindByID(id int) (*models.Book, error) {
	var gormBook gormmodels.GormBook

	err := r.db.First(&gormBook, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("kitap bulunamadı")
		}
		return nil, err
	}

	book := toBook(gormBook)
	return &book, nil
}

func (r *GormBookRepository) Create(book models.Book) (models.Book, error) {
	gormBook := toGormBook(book)

	err := r.db.Create(&gormBook).Error
	if err != nil {
		return models.Book{}, err
	}

	return toBook(gormBook), nil
}

func (r *GormBookRepository) Update(id int, book models.Book) (models.Book, error) {
	var gormBook gormmodels.GormBook

	err := r.db.First(&gormBook, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Book{}, errors.New("kitap bulunamadı")
		}
		return models.Book{}, err
	}

	gormBook.Title = book.Title
	gormBook.Author = book.Author
	gormBook.Year = book.Year

	err = r.db.Save(&gormBook).Error
	if err != nil {
		return models.Book{}, err
	}

	return toBook(gormBook), nil
}

func (r *GormBookRepository) Delete(id int) error {
	result := r.db.Delete(&gormmodels.GormBook{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("kitap bulunamadı")
	}

	return nil
}

func toBook(gormBook gormmodels.GormBook) models.Book {
	return models.Book{
		ID:     int(gormBook.ID),
		Title:  gormBook.Title,
		Author: gormBook.Author,
		Year:   gormBook.Year,
	}
}

func toGormBook(book models.Book) gormmodels.GormBook {
	return gormmodels.GormBook{
		ID:     uint(book.ID),
		Title:  book.Title,
		Author: book.Author,
		Year:   book.Year,
	}
}
