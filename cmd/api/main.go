package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/oopchoop/litebank/internal/auth"
	"github.com/oopchoop/litebank/internal/config"
	"github.com/oopchoop/litebank/internal/httpapi"
	"github.com/oopchoop/litebank/internal/ledger"
	"github.com/oopchoop/litebank/internal/platform/postgres"
	"github.com/oopchoop/litebank/internal/token"
	"github.com/rs/cors"

	_ "github.com/oopchoop/litebank/docs"
)

// @title Litebank API
// @version 1.0
// @description API for bank LiteBank
// @host localhost:8181
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.Load()

	if err != nil {
		log.Fatalf("get config error: %v", err)
	}

	pool, err := postgres.NewPool(context.Background(), cfg.DatabaseURL)

	if err != nil {
		log.Fatalf("db connection error: %v", err)
	}

	fmt.Println("success connect to db")

	defer pool.Close()

	JWTManager := token.NewJWTManager(cfg.JWTSecret)

	router := httpapi.NewRouter(httpapi.Handlers{
		DB:        pool,
		JWT:       JWTManager,
		AuthSvc:   auth.NewService(auth.NewRepository(pool), JWTManager),
		LedgerSvc: ledger.NewService(ledger.NewRepository(pool)),
	})

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Idempotency-Key"},
		AllowCredentials: true,
	})

	handler := httpapi.RecoveryMiddleware(c.Handler(router))

	server := &http.Server{
		Addr:              cfg.Host,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
