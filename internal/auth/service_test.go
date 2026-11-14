package auth

import (
	"context"
	"testing"
)

type mockUserRepo struct {
	savedUser NewUser
	mockError error
}

func (m *mockUserRepo) Create(ctx context.Context, user NewUser) (int64, error) {
	m.savedUser = user
	if m.mockError != nil {
		return 0, m.mockError
	}
	return 1, nil
}

func (m *mockUserRepo) GetByPhone(ctx context.Context, phone string) (*User, error) {
	return nil, nil
}

func TestService_Register(t *testing.T) {
	repo := &mockUserRepo{}

	svc := NewService(repo, nil)

	input := RegisterInput{
		PhoneNumber: "88881231251",
		Password:    "secret123321",
		FirstName:   "Ivan",
		LastName:    "Ivanov",
	}

	id, err := svc.Register(context.Background(), input)

	if err != nil {
		t.Fatalf("expected no err, get: %v", err)
	}

	if id != 1 {
		t.Fatalf("expected id 1, get: %d", id)
	}

	if repo.savedUser.PhoneNumber != "+78881231251" {
		t.Errorf("phone is not normalize, expected +78881231251, get: %s", repo.savedUser.PhoneNumber)
	}

	if repo.savedUser.PasswordHash == "secret123321" {
		t.Errorf("password saved open view")
	}

	if len(repo.savedUser.PasswordHash) < 50 {
		t.Errorf("password hash is too short, maybe bcrypt not working")
	}
}
