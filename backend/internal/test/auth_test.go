package internal_test

import (
	"context"
	"testing"

	"github.com/soydoradesu/product_discovery/internal/domain"
	"github.com/soydoradesu/product_discovery/internal/service"
)

type fakeUsers struct {
	byEmail map[string]domain.User
	byID    map[int64]domain.User
}

func (f *fakeUsers) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	u, ok := f.byEmail[email]
	if !ok {
		return domain.User{}, service.ErrUserNotFound
	}
	return u, nil
}

func (f *fakeUsers) GetByID(ctx context.Context, id int64) (domain.User, error) {
	u, ok := f.byID[id]
	if !ok {
		return domain.User{}, service.ErrUserNotFound
	}
	return u, nil
}

func TestAuthLogin_Success(t *testing.T) {
	hash, err := service.HashPassword("Password123!")
	if err != nil {
		t.Fatal(err)
	}

	f := &fakeUsers{
		byEmail: map[string]domain.User{
			"demo@example.com": {ID: 1, Email: "demo@example.com", PasswordHash: &hash},
		},
		byID: map[int64]domain.User{
			1: {ID: 1, Email: "demo@example.com", PasswordHash: &hash},
		},
	}

	svc := service.NewAuthService(f)

	id, err := svc.Login(context.Background(), "demo@example.com", "Password123!")
	if err != nil {
		t.Fatalf("expected ok, got err=%v", err)
	}
	if id != 1 {
		t.Fatalf("expected id=1 got %d", id)
	}
}

func TestAuthLogin_BadPassword(t *testing.T) {
	hash, _ := service.HashPassword("Password123!")
	f := &fakeUsers{
		byEmail: map[string]domain.User{
			"demo@example.com": {ID: 1, Email: "demo@example.com", PasswordHash: &hash},
		},
		byID: map[int64]domain.User{},
	}
	svc := service.NewAuthService(f)

	_, err := svc.Login(context.Background(), "demo@example.com", "wrong")
	if err == nil {
		t.Fatal("expected error")
	}
}
