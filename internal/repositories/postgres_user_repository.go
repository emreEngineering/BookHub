package repositories

import (
	"BookHub/internal/models"
	"database/sql"
	"errors"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

func (r *PostgresUserRepository) FindAll() ([]models.User, error) {
	rows, err := r.db.Query(`
		SELECT id, name, email, password_hash, role
		FROM users
		ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}

	for rows.Next() {
		var user models.User

		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *PostgresUserRepository) FindByID(id int) (*models.User, error) {
	var user models.User

	err := r.db.QueryRow(`
		SELECT id, name, email, password_hash, role
		FROM users
		WHERE id = $1
	`, id).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User

	err := r.db.QueryRow(`
		SELECT id, name, email, password_hash, role
		FROM users
		WHERE email = $1
	`, email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Role)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) Create(user models.User) (models.User, error) {
	if user.Role == "" {
		user.Role = "user"
	}

	err := r.db.QueryRow(`
		INSERT INTO users (name, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`, user.Name, user.Email, user.PasswordHash, user.Role).Scan(&user.ID)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
