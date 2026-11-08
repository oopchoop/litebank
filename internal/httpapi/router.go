package httpapi

import (
	"context"
	"net/http"

	"github.com/oopchoop/litebank/internal/auth"
	"github.com/oopchoop/litebank/internal/ledger"
	"github.com/oopchoop/litebank/internal/token"
)

type Pinger interface {
	Ping(ctx context.Context) error
}

type Handlers struct {
	DB        Pinger
	JWT       *token.JWTManager
	AuthSvc   *auth.Service
	LedgerSvc *ledger.Service
}

func NewRouter(h Handlers) http.Handler {
	mux := http.NewServeMux()

	// check db connection
	mux.HandleFunc("GET /health", HealthHandler(h.DB))

	// auth
	mux.HandleFunc("POST /register", RegisterHandler(h.AuthSvc))
	mux.HandleFunc("POST /auth", AuthHandler(h.AuthSvc))

	authMW := AuthMiddleware(h.JWT)

	// ledger
	mux.Handle("GET /accounts", authMW(GetAccountsHandler(h.LedgerSvc)))
	mux.Handle("GET /transfers", authMW(GetTransfersHandler(h.LedgerSvc)))

	mux.Handle("POST /accounts", authMW(CreateAccountHandler(h.LedgerSvc)))
	mux.Handle("POST /deposit", authMW(CreateDepositHandler(h.LedgerSvc)))
	mux.Handle("POST /transfer", authMW(CreateTransferHandler(h.LedgerSvc)))

	return RecoveryMiddleware(mux)
}
