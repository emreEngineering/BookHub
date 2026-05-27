package repositories

import "BookHub/internal/models"

type UserRepository interface {
	FindAll() ([]models.User, error)
	FindByID(id int) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Create(user models.User) (models.User, error)
}
