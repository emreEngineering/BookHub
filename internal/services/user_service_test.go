package services

import (
	"context"
	"testing"

	"BookHub/internal/models"
	"BookHub/internal/repositories"
)

type fakeUserRepository struct {
	users  []models.User
	nextID int
}

func newFakeUserRepository() *fakeUserRepository {
	return &fakeUserRepository{
		nextID: 1,
		users:  []models.User{},
	}
}

func (r *fakeUserRepository) FindAll(ctx context.Context) ([]models.User, error) {
	return r.users, nil
}

func (r *fakeUserRepository) FindByID(ctx context.Context, id int) (*models.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return &user, nil
		}
	}

	return nil, repositories.ErrUserNotFound
}

func (r *fakeUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	for _, user := range r.users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, repositories.ErrUserNotFound
}

func (r *fakeUserRepository) Create(ctx context.Context, user models.User) (models.User, error) {
	user.ID = r.nextID
	r.nextID++
	if user.Role == "" {
		user.Role = "user"
	}
	r.users = append(r.users, user)
	return user, nil
}

func TestUserService_Register_Success(t *testing.T) {
	service := NewUserService(newFakeUserRepository())

	user, err := service.Register(context.Background(), RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "123456",
	})
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	if user.ID == 0 {
		t.Fatal("expected created user ID to be set")
	}

	if user.PasswordHash == "123456" {
		t.Fatal("expected password hash not to equal plain password")
	}
}

func TestUserService_Register_ValidationError(t *testing.T) {
	service := NewUserService(newFakeUserRepository())

	_, err := service.Register(context.Background(), RegisterRequest{
		Name:     "",
		Email:    "test@example.com",
		Password: "123456",
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestUserService_Register_DuplicateEmail(t *testing.T) {
	service := NewUserService(newFakeUserRepository())
	request := RegisterRequest{
		Name:     "Test User",
		Email:    "duplicate@example.com",
		Password: "123456",
	}

	_, err := service.Register(context.Background(), request)
	if err != nil {
		t.Fatalf("first Register returned error: %v", err)
	}

	_, err = service.Register(context.Background(), request)
	if err == nil {
		t.Fatal("expected duplicate email error")
	}
}

func TestUserService_Login_Success(t *testing.T) {
	service := NewUserService(newFakeUserRepository())

	registeredUser, err := service.Register(context.Background(), RegisterRequest{
		Name:     "Login User",
		Email:    "login@example.com",
		Password: "123456",
	})
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	user, err := service.Login(context.Background(), LoginRequest{
		Email:    "login@example.com",
		Password: "123456",
	})
	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}

	if user.ID != registeredUser.ID {
		t.Fatalf("expected user ID %d, got %d", registeredUser.ID, user.ID)
	}
}

func TestUserService_Login_WrongPassword(t *testing.T) {
	service := NewUserService(newFakeUserRepository())

	_, err := service.Register(context.Background(), RegisterRequest{
		Name:     "Login User",
		Email:    "wrong-password@example.com",
		Password: "123456",
	})
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	_, err = service.Login(context.Background(), LoginRequest{
		Email:    "wrong-password@example.com",
		Password: "wrong",
	})
	if err == nil {
		t.Fatal("expected wrong password error")
	}
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	service := NewUserService(newFakeUserRepository())

	_, err := service.Login(context.Background(), LoginRequest{
		Email:    "missing@example.com",
		Password: "123456",
	})
	if err == nil {
		t.Fatal("expected missing user error")
	}
}

func TestUserService_GetUserByID_Found(t *testing.T) {
	service := NewUserService(newFakeUserRepository())

	registeredUser, err := service.Register(context.Background(), RegisterRequest{
		Name:     "Lookup User",
		Email:    "lookup@example.com",
		Password: "123456",
	})
	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	user, err := service.GetUserByID(context.Background(), registeredUser.ID)
	if err != nil {
		t.Fatalf("GetUserByID returned error: %v", err)
	}

	if user == nil || user.ID != registeredUser.ID {
		t.Fatalf("expected user ID %d, got %#v", registeredUser.ID, user)
	}
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	service := NewUserService(newFakeUserRepository())

	user, err := service.GetUserByID(context.Background(), 999)
	if err == nil {
		t.Fatal("expected missing user error")
	}

	if user != nil {
		t.Fatalf("expected nil user, got %#v", user)
	}
}
