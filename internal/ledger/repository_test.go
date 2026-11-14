package ledger

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("Not find .env, use system variablies")
	}

	dbURL := os.Getenv("DATABASE_TEST_URL")
	if dbURL == "" {
		t.Fatalf("DATABASE_TEST_URL not set in .env")
	}

	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}

	_, err = pool.Exec(context.Background(),
		`TRUNCATE accounts, users_accounts, transfers, entries
		 RESTART IDENTITY CASCADE`)
	if err != nil {
		t.Fatalf("failed truncate tables: %v", err)
	}

	return pool
}

func TestRepository_Transfer(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()

	repo := NewRepository(pool)
	ctx := context.Background()

	var userID int64
	err := pool.QueryRow(ctx,
		`INSERT INTO users (phone_number, password_hash, first_name, last_name)
		 VALUES('+79990000000', 'hash', 'Ivan', 'Ivanov')
		 RETURNING id`).Scan(&userID)
	if err != nil {
		t.Fatal(err)
	}

	var accFromID, accToID int64

	err = pool.QueryRow(ctx,
		`INSERT INTO accounts (account_type_id, balance)
		 VALUES (2, 10000)
		 RETURNING id`).Scan(&accFromID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = pool.Exec(ctx,
		`INSERT INTO users_accounts (user_id, account_id)
		 VALUES($1, $2)`, userID, accFromID)
	if err != nil {
		t.Fatal(err)
	}

	err = pool.QueryRow(ctx,
		`INSERT INTO accounts (account_type_id, balance)
		 VALUES(2, 0)
		 RETURNING id`).Scan(&accToID)
	if err != nil {
		t.Fatal(err)
	}

	err = repo.Transfer(ctx, NewTransfer{
		UserID:        userID,
		AccountFromID: accFromID,
		AccountToID:   accToID,
		Amount:        5000,
		IsSystem:      false,
	})
	if err != nil {
		t.Fatal(err)
	}

	var balanceFrom, balanceTo int64
	pool.QueryRow(ctx, `SELECT balance FROM accounts WHERE id = $1`, accFromID).Scan(&balanceFrom)
	pool.QueryRow(ctx, `SELECT balance FROM accounts WHERE id = $1`, accToID).Scan(&balanceTo)

	if balanceFrom != 5000 {
		t.Errorf("expected from balance 5000, get: %d", balanceFrom)
	}

	if balanceTo != 5000 {
		t.Errorf("expected to balance 5000, get: %d", balanceTo)
	}

	var entriesCount int64
	pool.QueryRow(ctx, `SELECT COUNT(*) FROM entries`).Scan(&entriesCount)
	if entriesCount != 2 {
		t.Errorf("expected 2 entries, get: %d", entriesCount)
	}
}
