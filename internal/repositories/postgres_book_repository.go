package repositories

import (
	"BookHub/internal/models"
	"database/sql"
	"errors"
)

type PostgresBookRepository struct {
	db *sql.DB
}

func NewPostgresBookRepository(db *sql.DB) *PostgresBookRepository {
	return &PostgresBookRepository{
		db: db,
	}
}

func (r *PostgresBookRepository) FindAll() ([]models.Book, error) {
	rows, err := r.db.Query(`
		SELECT id, title, author, year
		FROM books
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := []models.Book{}

	for rows.Next() {
		var book models.Book

		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (r *PostgresBookRepository) FindByID(id int) (*models.Book, error) {
	var book models.Book

	err := r.db.QueryRow(`
		SELECT id, title, author, year
		FROM books
		WHERE id = $1
	`, id).Scan(&book.ID, &book.Title, &book.Author, &book.Year)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("kitap bulunamadı")
		}
		return nil, err
	}
	return &book, nil
}

func (r *PostgresBookRepository) Create(book models.Book) (models.Book, error) {
	err := r.db.QueryRow(`
		INSERT INTO books (title, author, year)
		VALUES ($1, $2, $3)
		RETURNING id
	`, book.Title, book.Author, book.Year).Scan(&book.ID)

	if err != nil {
		return models.Book{}, err
	}

	return book, nil
}

func (r *PostgresBookRepository) Update(id int, book models.Book) (models.Book, error) {
	result, err := r.db.Exec(`
		UPDATE books
		SET title = $1, author = $2, year = $3
		WHERE id = $4
	`, book.Title, book.Author, book.Year, id)

	if err != nil {
		return models.Book{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Book{}, err
	}

	if rowsAffected == 0 {
		return models.Book{}, errors.New("kitap bulunamadı")
	}

	book.ID = id

	return book, nil
}

func (r *PostgresBookRepository) Delete(id int) error {
	result, err := r.db.Exec(`
		DELETE FROM books
		WHERE id = $1
	`, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("kitap bulunamadı")
	}

	return nil
}
