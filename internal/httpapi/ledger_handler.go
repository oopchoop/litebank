package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/oopchoop/litebank/internal/ledger"
)

type CreateAccountResponse struct {
	ID int64 `json:"id"`
}

type CreateDepositRequest struct {
	AccountToID int64 `json:"account_to_id"`
	Amount      int64 `json:"amount"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type CreateTransferRequest struct {
	AccountToID   int64 `json:"account_to_id"`
	AccountFromID int64 `json:"account_from_id"`
	Amount        int64 `json:"amount"`
}

func CreateAccountHandler(service *ledger.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserID(r.Context())

		if !ok {
			writeError(w, http.StatusUnauthorized, "unauthorized")

			return
		}

		accountID, err := service.CreateAccount(r.Context(), userID)

		if err != nil {
			writeError(w, statusForError(err), publicMessage(err))
			return
		}

		writeJSON(w, http.StatusCreated, CreateAccountResponse{
			ID: accountID,
		})
	}
}

func CreateDepositHandler(service *ledger.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		userID, ok := GetUserID(ctx)

		if !ok {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		var req CreateDepositRequest
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		if err := dec.Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		idempotencyKey := r.Header.Get("Idempotency-Key")

		err := service.CreateDeposit(ctx, ledger.NewDeposit{
			UserID:         userID,
			AccountToID:    req.AccountToID,
			Amount:         req.Amount,
			IdempotencyKey: idempotencyKey,
		})

		if err != nil {
			writeError(w, statusForError(err), publicMessage(err))
			return
		}

		writeJSON(w, http.StatusOK, MessageResponse{
			Message: "success deposit",
		})
	}
}

func CreateTransferHandler(service *ledger.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)

		defer cancel()

		userID, ok := GetUserID(ctx)

		if !ok {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		var req CreateTransferRequest
		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()

		if err := dec.Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		idempotencyKey := r.Header.Get("Idempotency-Key")

		err := service.CreateTransfer(ctx, ledger.NewTransfer{
			UserID:         userID,
			AccountToID:    req.AccountToID,
			AccountFromID:  req.AccountFromID,
			Amount:         req.Amount,
			IsSystem:       false,
			IdempotencyKey: idempotencyKey,
		})

		if err != nil {
			writeError(w, statusForError(err), publicMessage(err))
			return
		}

		writeJSON(w, http.StatusCreated, MessageResponse{
			Message: "success create transfer",
		})
	}
}

func GetAccountsHandler(service *ledger.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)

		defer cancel()

		userID, ok := GetUserID(ctx)

		if !ok {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		accounts, err := service.GetAccounts(ctx, userID)

		if err != nil {
			writeError(w, statusForError(err), publicMessage(err))
			return
		}

		writeJSON(w, http.StatusOK, accounts)
	}
}

func GetTransfersHandler(service *ledger.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)

		defer cancel()

		userID, ok := GetUserID(ctx)

		if !ok {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		transfers, err := service.GetTransfers(ctx, userID)

		if err != nil {
			writeError(w, statusForError(err), publicMessage(err))
			return
		}

		writeJSON(w, http.StatusOK, transfers)
	}
}
