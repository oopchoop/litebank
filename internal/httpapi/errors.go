package httpapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/oopchoop/litebank/internal/auth"
)

func statusForError(err error) int {
	switch {
	case errors.Is(err, auth.ErrValidation):
		return http.StatusBadRequest

	case errors.Is(err, auth.ErrPhoneAlreadyUsed):
		return http.StatusConflict

	case errors.Is(err, context.DeadlineExceeded):
		return http.StatusGatewayTimeout

	case errors.Is(err, context.Canceled):
		return 499

	default:
		return http.StatusInternalServerError
	}
}

func publicMessage(err error) string {
	switch statusForError(err) {
	case http.StatusBadRequest, http.StatusConflict:
		return err.Error()
	case http.StatusGatewayTimeout:
		return "request timeout"
	case 499:
		return "request canceled"
	default:
		return "internal server error"
	}
}
