package auth

import (
	"strings"
	"unicode/utf8"
)

const maxNameLen = 255

func normalizePhone(raw string) (string, error) {
	var digits strings.Builder

	for _, r := range raw {
		if r >= '0' && r <= '9' {
			digits.WriteRune(r)
		}
	}

	d := digits.String()

	switch len(d) {
	case 10:
		if d[0] != '9' {
			return "", ValidationError(ErrIncorrectPhone)
		}

		d = "7" + d
	case 11:
		if d[0] == '8' {
			d = "7" + d[1:]
		} else if d[0] != '7' {
			return "", ValidationError(ErrIncorrectPhone)
		}

	default:
		return "", ValidationError(ErrIncorrectPhone)
	}

	return "+" + d, nil
}

func validatePassword(raw string) error {
	if utf8.RuneCountInString(raw) < 8 {
		return ValidationError(ErrPasswordTooShort)
	}

	if len(raw) > 72 {
		return ValidationError(ErrPasswordTooLong)
	}

	return nil
}

func normalizeNames(first, last string) (string, string, error) {
	trimmedFirst := strings.TrimSpace(first)
	trimmedLast := strings.TrimSpace(last)

	if trimmedFirst == "" || trimmedLast == "" {
		return "", "", ValidationError(ErrNameEmpty)
	}

	if utf8.RuneCountInString(trimmedFirst) > maxNameLen || utf8.RuneCountInString(trimmedLast) > maxNameLen {
		return "", "", ValidationError(ErrNameTooLong)
	}

	return trimmedFirst, trimmedLast, nil
}

func normalizePatronymic(raw string) (*string, error) {
	trimmedRaw := strings.TrimSpace(raw)

	if trimmedRaw == "" {
		return nil, nil
	}

	if utf8.RuneCountInString(trimmedRaw) > maxNameLen {
		return nil, ValidationError(ErrPatronymicTooLong)
	}

	return &trimmedRaw, nil
}
