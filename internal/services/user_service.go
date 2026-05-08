package services

import (
	"errors"

	"BookHub/internal/models"
	"BookHub/internal/repositories"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserService interface {
	Register(request RegisterRequest) (models.User, error)
	GetAllUsers() ([]models.User, error)
	GetUserByID(id int) (*models.User, error)
}

type DefaultUserService struct {
	userRepo repositories.UserRepositories
}

func NewUserService(userRepo repositories.UserRepositories) *DefaultUserService {
	return &DefaultUserService{
		userRepo: userRepo,
	}
}

func (s *DefaultUserService) Register(request RegisterRequest) (models.User, error) {
	err := validateRegisterRequest(request)
	if err != nil {
		return models.User{}, err
	}

	_, err = s.userRepo.FindByEmail(request.Email)
	if err == nil {
		return models.User{}, errors.New("bu email zaten kayıtlı")
	}

	user := models.User{
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: request.Password,
		Role:         "user",
	}

	return s.userRepo.Create(user)
}

func (s *DefaultUserService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.FindAll()
}

func (s *DefaultUserService) GetUserByID(id int) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func validateRegisterRequest(request RegisterRequest) error {
	if request.Name == "" {
		return errors.New("isim boş olamaz")
	}

	if request.Email == "" {
		return errors.New("email boş olamaz")
	}

	if request.Password == "" {
		return errors.New("şifre boş olamaz")
	}

	return nil
}
