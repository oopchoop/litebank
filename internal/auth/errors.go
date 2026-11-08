package auth

import (
	"errors"
	"fmt"
)

var (
	// marker validation error
	ErrValidation = errors.New("validation error")

	// phone errors
	ErrPhoneAlreadyUsed = errors.New("phone already used")
	ErrIncorrectPhone   = errors.New("incorrect phone format")
	ErrUserNotRegister  = errors.New("user not register")

	//password errors
	ErrPasswordTooShort = errors.New("password too short")
	ErrPasswordTooLong  = errors.New("password too long")

	// name errors
	ErrNameTooLong = errors.New("last or first name too long")
	ErrNameEmpty   = errors.New("last or first name is empty")

	// patronymic errors
	ErrPatronymicTooLong = errors.New("patronymic too long")

	// auth data errors
	ErrIncorrectPasswordOrPhone = errors.New("incorrect phone or password")
)

func ValidationError(err error) error {
	return fmt.Errorf("%w: %w", ErrValidation, err)
}
