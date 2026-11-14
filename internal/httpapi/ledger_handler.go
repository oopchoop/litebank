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

// CreateAccountHandler godoc
// @Summary Create account
// @Description Create new account for user
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 201 {object} CreateAccountResponse "Success created account"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /accounts [post]
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

// CreateDepositHandler godoc
// @Summary Create deposit
// @Description Create deposit for account
// @Tags transfers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Idempotency-Key header string false "Unique key to prevent double deposit"
// @Param input body CreateDepositRequest true "Data for create deposit"
// @Success 200 {object} MessageResponse "Success deposit"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Access denied"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /deposit [post]
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

// CreateTransferHandler godoc
// @Summary Create transfer
// @Description Create transfer from account to another account
// @Tags transfers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Idempotency-Key header string false "Unique key to prevent double tranfser"
// @Param input body CreateTransferRequest true "Data for create tranfser"
// @Success 201 {object} MessageResponse "Success create transfer"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 403 {object} ErrorResponse "Access denied"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /transfer [post]
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

// GetAccountHandler godoc
// @Summary Get user accounts
// @Description Get accounts for user
// @Tags accounts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} ledger.Account
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /accounts [get]
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

// GetTransfersHandler godoc
// @Summary Get user transfers
// @Description Get transfers for user
// @Tags transfers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} ledger.TransfersList
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /transfers [get]
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
