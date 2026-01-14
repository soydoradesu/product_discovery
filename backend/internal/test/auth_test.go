package internal_test

import (
	"context"
	"testing"

	"github.com/soydoradesu/product_discovery/internal/domain"
	"github.com/soydoradesu/product_discovery/internal/repository"
	"github.com/soydoradesu/product_discovery/internal/service"
)

type fakeUsers struct {
	byEmail map[string]domain.User
	byID    map[int64]domain.User
}

func (f *fakeUsers) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	u, ok := f.byEmail[email]
	if !ok {
		return domain.User{}, repository.ErrNotFound
	}
	return u, nil
}

func (f *fakeUsers) GetByID(ctx context.Context, id int64) (domain.User, error) {
	u, ok := f.byID[id]
	if !ok {
		return domain.User{}, repository.ErrNotFound
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

func (f *fakeUsers) GetByGoogleID(ctx context.Context, googleID string) (domain.User, error) {
	for _, u := range f.byEmail {
		if u.GoogleID != nil && *u.GoogleID == googleID {
			return u, nil
		}
	}
	return domain.User{}, repository.ErrNotFound
}

func (f *fakeUsers) SetGoogleID(ctx context.Context, userID int64, googleID string) error {
	u, ok := f.byID[userID]
	if !ok {
		return repository.ErrNotFound
	}
	u.GoogleID = &googleID
	f.byID[userID] = u
	// keep byEmail in sync
	if u.Email != "" {
		f.byEmail[u.Email] = u
	}
	return nil
}

func (f *fakeUsers) CreateOAuthUser(ctx context.Context, email, googleID string) (int64, error) {
	id := int64(len(f.byID) + 1)
	u := domain.User{ID: id, Email: email, GoogleID: &googleID}
	if f.byEmail == nil {
		f.byEmail = map[string]domain.User{}
	}
	if f.byID == nil {
		f.byID = map[int64]domain.User{}
	}
	f.byEmail[email] = u
	f.byID[id] = u
	return id, nil
}