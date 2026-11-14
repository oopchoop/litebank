package auth

import (
	"errors"
	"strings"
	"testing"
)

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name          string
		password      string
		expectedError error
	}{
		{
			name:          "short password",
			password:      "123",
			expectedError: ErrPasswordTooShort,
		},
		{
			name:          "valid password",
			password:      "supersecret",
			expectedError: nil,
		},
		{
			name:          "long password",
			password:      strings.Repeat("a", 100),
			expectedError: ErrPasswordTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePassword(tt.password)

			if !errors.Is(err, tt.expectedError) {
				t.Errorf("validatePassword(%q) = %v; expected %v", tt.password, err, tt.expectedError)
			}
		})
	}
}

func TestNormalizePhone(t *testing.T) {
	tests := []struct {
		name          string
		phone         string
		expectedPhone string
		expectedError error
	}{
		{
			name:          "short phone",
			phone:         "123",
			expectedPhone: "",
			expectedError: ErrIncorrectPhone,
		},
		{
			name:          "start 8",
			phone:         "88005553535",
			expectedPhone: "+78005553535",
			expectedError: nil,
		},
		{
			name:          "start 7",
			phone:         "78005553535",
			expectedPhone: "+78005553535",
			expectedError: nil,
		},
		{
			name:          "10 number",
			phone:         "9991234567",
			expectedPhone: "+79991234567",
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			phone, err := normalizePhone(tt.phone)

			if !errors.Is(err, tt.expectedError) {
				t.Errorf("normalizePhone(%q) = %v; expected %v", tt.phone, err, tt.expectedError)
			}

			if phone != tt.expectedPhone {
				t.Errorf("normalizePhone(%q) result = %v; currentResult = %v", tt.phone, phone, tt.expectedPhone)
			}
		})
	}
}
