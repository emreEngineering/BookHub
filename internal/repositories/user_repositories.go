package repositories

import (
	"BookHub/internal/models"
	"errors"
)

type UserRepositories interface {
	FindAll() ([]models.User, error)
	FindByID(id int) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Create(user models.User) (*models.User, error)
}

type MemoryUserRepository struct {
	users  []models.User
	nextID int
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		nextID: 2,
		users: []models.User{
			{
				ID:           1,
				Name:         "Admin User",
				Email:        "admin@bookhub.com",
				PasswordHash: "admin123",
				Role:         "admin",
			},
		},
	}
}
func (r *MemoryUserRepository) FindAll() ([]models.User, error) {
	return r.users, nil
}

func (r *MemoryUserRepository) FindByID(id int) (*models.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return &user, nil
		}
	}

	return nil, errors.New("kullanıcı bulunamadı")
}

func (r *MemoryUserRepository) FindByEmail(email string) (*models.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return &user, nil
		}
	}
	return nil, errors.New("kullanıcı bulunamadı")
}

func (r *MemoryUserRepository) Create(user models.User) (models.User, error) {
	user.ID = r.nextID
	r.nextID++

	if user.Role == "" {
		user.Role = "user"
	}
	r.users = append(r.users, user)
	return user, nil
}
