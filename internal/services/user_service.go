package services

import (
	"context"
	"errors"

	"BookHub/internal/activity"
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
	Register(ctx context.Context, request RegisterRequest) (models.User, error)
	Login(ctx context.Context, request LoginRequest) (models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
}

type DefaultUserService struct {
	userRepo       repositories.UserRepository
	activityLogger activity.ActivityLogger
}

func NewUserService(userRepo repositories.UserRepository, activityLoggers ...activity.ActivityLogger) *DefaultUserService {
	service := &DefaultUserService{
		userRepo: userRepo,
	}

	if len(activityLoggers) > 0 {
		service.activityLogger = activityLoggers[0]
	}

	return service
}

func (s *DefaultUserService) Register(ctx context.Context, request RegisterRequest) (models.User, error) {
	err := validateRegisterRequest(request)
	if err != nil {
		return models.User{}, err
	}

	_, err = s.userRepo.FindByEmail(ctx, request.Email)
	if err == nil {
		return models.User{}, repositories.ErrDuplicateEmail
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

	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		if errors.Is(err, repositories.ErrDuplicateEmail) {
			return models.User{}, repositories.ErrDuplicateEmail
		}
		return models.User{}, err
	}

	s.logActivity(ctx, "user_registered", "Kullanıcı kayıt oldu", &createdUser.ID, map[string]interface{}{
		"email": createdUser.Email,
		"name":  createdUser.Name,
		"role":  createdUser.Role,
	})

	return createdUser, nil
}

func (s *DefaultUserService) Login(ctx context.Context, request LoginRequest) (models.User, error) {
	err := validateLoginRequest(request)

	if err != nil {
		return models.User{}, err
	}
	user, err := s.userRepo.FindByEmail(ctx, request.Email)
	if err != nil {
		return models.User{}, errors.New("email veya şifre hatalı")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))
	if err != nil {
		return models.User{}, errors.New("email veya şifre hatalı")
	}

	s.logActivity(ctx, "user_login", "Kullanıcı giriş yaptı", &user.ID, map[string]interface{}{
		"email": user.Email,
	})

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

func (s *DefaultUserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepo.FindAll(ctx)
}

func (s *DefaultUserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	return s.userRepo.FindByID(ctx, id)
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

func (s *DefaultUserService) logActivity(ctx context.Context, eventType string, message string, userID *int, metadata map[string]interface{}) {
	if s.activityLogger == nil {
		return
	}

	_ = s.activityLogger.Log(ctx, eventType, message, userID, metadata)
}
