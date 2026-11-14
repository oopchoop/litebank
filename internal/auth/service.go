package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/oopchoop/litebank/internal/token"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, user NewUser) (int64, error)
	GetByPhone(ctx context.Context, phone string) (*User, error)
}

type Service struct {
	repo       UserRepository
	tokenMaker *token.JWTManager
}

func NewService(repo UserRepository, tokenMaker *token.JWTManager) *Service {
	return &Service{
		repo:       repo,
		tokenMaker: tokenMaker,
	}
}

func (s *Service) Register(ctx context.Context, in RegisterInput) (int64, error) {
	normalizedPhone, err := normalizePhone(in.PhoneNumber)

	if err != nil {
		return 0, err
	}

	if err := validatePassword(in.Password); err != nil {
		return 0, err
	}

	firstName, lastName, err := normalizeNames(in.FirstName, in.LastName)

	if err != nil {
		return 0, err
	}

	patronymic, err := normalizePatronymic(in.Patronymic)

	if err != nil {
		return 0, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)

	if err != nil {
		return 0, fmt.Errorf("hash password: %w", err)
	}

	return s.repo.Create(ctx, NewUser{
		PhoneNumber:  normalizedPhone,
		PasswordHash: string(hash),
		FirstName:    firstName,
		LastName:     lastName,
		Patronymic:   patronymic,
	})
}

func (s *Service) Login(ctx context.Context, in AuthInput) (string, error) {
	normalizePhone, err := normalizePhone(in.PhoneNumber)

	if err != nil {
		return "", err
	}

	user, err := s.repo.GetByPhone(ctx, normalizePhone)

	if err != nil {
		if errors.Is(err, ErrUserNotRegister) {
			return "", ValidationError(ErrIncorrectPasswordOrPhone)
		}

		return "", err
	}

	correctPassword := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password))

	if correctPassword != nil {
		return "", ErrIncorrectPasswordOrPhone
	}

	accessToken, err := s.tokenMaker.Generate(user.ID)

	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	return accessToken, nil
}
