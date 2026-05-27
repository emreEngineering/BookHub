package services

import (
	"errors"

	"BookHub/internal/models"
	"BookHub/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserService interface {
	Register(request RegisterRequest) (models.User, error)
	Login(request LoginRequest) (models.User, error)
	GetAllUsers() ([]models.User, error)
	GetUserByID(id int) (*models.User, error)
}

type DefaultUserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *DefaultUserService {
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
	if !errors.Is(err, repositories.ErrUserNotFound) {
		return models.User{}, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: string(passwordHash),
		Role:         "user",
	}

	return s.userRepo.Create(user)
}

func (s *DefaultUserService) Login(request LoginRequest) (models.User, error) {
	err := validateLoginRequest(request)

	if err != nil {
		return models.User{}, err
	}
	user, err := s.userRepo.FindByEmail(request.Email)
	if err != nil {
		return models.User{}, errors.New("email veya şifre hatalı")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		return models.User{}, errors.New("email veya şifre hatalı")
	}
	return *user, nil
}

func validateLoginRequest(request LoginRequest) error {
	if request.Email == "" {
		return errors.New("email boş olamaz")
	}
	if request.Password == "" {
		return errors.New("şifre boş olamaz")
	}
	return nil
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
