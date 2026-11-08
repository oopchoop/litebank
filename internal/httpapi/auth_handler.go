package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/oopchoop/litebank/internal/auth"
)

type RegisterRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Patronymic  string `json:"patronymic"`
}

type RegisterResponse struct {
	ID int64 `json:"id"`
}

type AuthRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type AuthResponse struct {
	UserToken string `json:"user_token"`
}

func RegisterHandler(service *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		var req RegisterRequest
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		if err := dec.Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		id, err := service.Register(ctx, auth.RegisterInput{
			PhoneNumber: req.PhoneNumber,
			Password:    req.Password,
			FirstName:   req.FirstName,
			LastName:    req.LastName,
			Patronymic:  req.Patronymic,
		})

		if err != nil {
			writeError(w, statusForError(err), publicMessage(err))
			return
		}

		writeJSON(w, http.StatusCreated, RegisterResponse{ID: id})
	}
}

func AuthHandler(service *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		var req AuthRequest
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		if err := dec.Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		token, err := service.Login(ctx, auth.AuthInput{
			PhoneNumber: req.PhoneNumber,
			Password:    req.Password,
		})

		if err != nil {
			writeError(w, statusForError(err), publicMessage(err))
			return
		}

		writeJSON(w, http.StatusOK, AuthResponse{
			UserToken: token,
		})
	}
}
