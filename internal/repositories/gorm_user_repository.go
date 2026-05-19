package repositories

import (
	gormmodels "BookHub/internal/gormdb/models"
	"BookHub/internal/models"
	"errors"

	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{
		db: db,
	}
}

func (r *GormUserRepository) FindAll() ([]models.User, error) {
	var gormUsers []gormmodels.GormUser

	err := r.db.Order("id").Find(&gormUsers).Error
	if err != nil {
		return nil, err
	}

	users := make([]models.User, 0, len(gormUsers))
	for _, gormUser := range gormUsers {
		users = append(users, toUser(gormUser))
	}

	return users, nil
}

func (r *GormUserRepository) FindByID(id int) (*models.User, error) {
	var gormUser gormmodels.GormUser

	err := r.db.First(&gormUser, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("kullanıcı bulunamadı")
		}
		return nil, err
	}

	user := toUser(gormUser)
	return &user, nil
}

func (r *GormUserRepository) FindByEmail(email string) (*models.User, error) {
	var gormUser gormmodels.GormUser

	err := r.db.Where("email = ?", email).First(&gormUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("kullanıcı bulunamadı")
		}
		return nil, err
	}

	user := toUser(gormUser)
	return &user, nil
}

func (r *GormUserRepository) Create(user models.User) (models.User, error) {
	if user.Role == "" {
		user.Role = "user"
	}

	gormUser := toGormUser(user)

	err := r.db.Create(&gormUser).Error
	if err != nil {
		return models.User{}, err
	}

	return toUser(gormUser), nil
}

func toUser(gormUser gormmodels.GormUser) models.User {
	return models.User{
		ID:           int(gormUser.ID),
		Name:         gormUser.Name,
		Email:        gormUser.Email,
		PasswordHash: gormUser.PasswordHash,
		Role:         gormUser.Role,
	}
}

func toGormUser(user models.User) gormmodels.GormUser {
	return gormmodels.GormUser{
		ID:           uint(user.ID),
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
	}
}
