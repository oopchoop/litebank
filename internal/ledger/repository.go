package ledger

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateAccount(ctx context.Context, userID int64) (int64, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}

	defer tx.Rollback(ctx)

	var id int64

	errInsert := tx.QueryRow(ctx,
		`INSERT INTO accounts(account_type_id)
			SELECT id 
			FROM accounts_types 
			WHERE name = 'personal'
		RETURNING id`,
	).Scan(&id)

	if errInsert != nil {
		return 0, fmt.Errorf("account create: %w", errInsert)
	}

	_, linkErr := tx.Exec(ctx,
		`INSERT INTO users_accounts(user_id, account_id) 
		 VALUES($1, $2)`,
		userID, id)

	if linkErr != nil {
		return 0, fmt.Errorf("link account with user: %w", linkErr)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit tx: %w", err)
	}

	return id, nil
}

func (r *Repository) executeTransfer(ctx context.Context, transfer NewTransfer, tx pgx.Tx) error {
	_, err := tx.Exec(ctx,
		`SELECT 1
		 FROM accounts
		 WHERE id IN ($1, $2)
		 ORDER BY id
		 FOR UPDATE`,
		transfer.AccountFromID, transfer.AccountToID,
	)

	if err != nil {
		return fmt.Errorf("select accounts: %w", err)
	}

	if !transfer.IsSystem {
		if err := r.checkSufficientFunds(ctx, transfer.AccountFromID, transfer.Amount, tx); err != nil {
			return err
		}
	}

	var transferID int64

	transferErr := tx.QueryRow(ctx,
		`INSERT INTO transfers(amount, account_from_id, account_to_id, idempotency_key)
		VALUES ($1, $2, $3, $4)
		RETURNING id`,
		transfer.Amount, transfer.AccountFromID, transfer.AccountToID, transfer.IdempotencyKey,
	).Scan(&transferID)

	if transferErr != nil {
		return fmt.Errorf("create transfer: %w", transferErr)
	}

	_, entriesErr := tx.Exec(ctx,
		`INSERT INTO entries(transfer_id, account_id, amount)
		 VALUES ($1, $2, $3), ($1, $4, $5)`,
		transferID, transfer.AccountToID, transfer.Amount, transfer.AccountFromID, -transfer.Amount,
	)

	if entriesErr != nil {
		return fmt.Errorf("create entries: %w", entriesErr)
	}

	_, updateSystemAccountErr := tx.Exec(ctx,
		`UPDATE accounts
		 SET balance = balance - $1
		 WHERE id = $2`,
		transfer.Amount, transfer.AccountFromID,
	)

	if updateSystemAccountErr != nil {
		return fmt.Errorf("update from account: %w", updateSystemAccountErr)
	}

	_, updateAccountErr := tx.Exec(ctx,
		`UPDATE accounts
		 SET balance = balance + $1
		 WHERE id = $2`,
		transfer.Amount, transfer.AccountToID,
	)

	if updateAccountErr != nil {
		return fmt.Errorf("update account: %w", updateAccountErr)
	}

	return nil
}

func (r *Repository) checkAccountOwner(ctx context.Context, userID int64, accountID int64, tx pgx.Tx) error {
	var isOwnerAccount int64

	isOwnerAccountErr := tx.QueryRow(ctx,
		`SELECT 1
		 FROM users_accounts
		 WHERE user_id = $1 AND account_id = $2`,
		userID, accountID,
	).Scan(&isOwnerAccount)

	if isOwnerAccountErr != nil {
		if errors.Is(isOwnerAccountErr, pgx.ErrNoRows) {
			return ErrAccessDenied
		}

		return fmt.Errorf("check account owner: %w", isOwnerAccountErr)
	}

	return nil
}

func (r *Repository) checkSufficientFunds(ctx context.Context, accountID int64, transferAmount int64, tx pgx.Tx) error {
	var balanceCurrentAccount int64

	balanceCurrentAccountErr := tx.QueryRow(ctx,
		`SELECT balance
		 FROM accounts
		 WHERE id = $1
		 FOR UPDATE`,
		accountID,
	).Scan(&balanceCurrentAccount)

	if balanceCurrentAccountErr != nil {
		return fmt.Errorf("internal server: %w", balanceCurrentAccountErr)
	}

	if balanceCurrentAccount-transferAmount < 0 {
		return ErrNoBalanceForTransfer
	}

	return nil
}

func (r *Repository) beginTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}

	return tx, nil
}

func (r *Repository) checkIdempotencyKey(ctx context.Context, idempotencyKey string) (bool, error) {
	if idempotencyKey == "" {
		return false, nil
	}

	var exists bool

	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM transfers WHERE idempotency_key = $1)`,
		idempotencyKey,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("Select exists transfer for idempotency_key: %w", err)
	}

	return exists, nil
}

func (r *Repository) Transfer(ctx context.Context, transfer NewTransfer) error {
	exists, err := r.checkIdempotencyKey(ctx, transfer.IdempotencyKey)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	tx, err := r.beginTx(ctx)

	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	if err := r.checkAccountOwner(ctx, transfer.UserID, transfer.AccountFromID, tx); err != nil {
		return err
	}

	if err := r.executeTransfer(ctx, transfer, tx); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil
		}

		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (r *Repository) Deposit(ctx context.Context, deposit NewDeposit) error {
	exists, err := r.checkIdempotencyKey(ctx, deposit.IdempotencyKey)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	tx, err := r.beginTx(ctx)

	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	if err := r.checkAccountOwner(ctx, deposit.UserID, deposit.AccountToID, tx); err != nil {
		return err
	}

	var systemAccountID int64

	systemAccountErr := tx.QueryRow(ctx,
		`SELECT id
		 FROM accounts
		 WHERE is_system = true`,
	).Scan(&systemAccountID)

	if systemAccountErr != nil {
		return fmt.Errorf("system account: %w", systemAccountErr)
	}

	if err := r.executeTransfer(ctx, NewTransfer{
		UserID:         deposit.UserID,
		AccountToID:    deposit.AccountToID,
		AccountFromID:  systemAccountID,
		Amount:         deposit.Amount,
		IsSystem:       true,
		IdempotencyKey: deposit.IdempotencyKey,
	}, tx); err != nil {
		var pgxErr *pgconn.PgError
		if errors.As(err, &pgxErr) && pgxErr.Code == "23505" {
			return nil
		}

		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (r *Repository) GetAccounts(ctx context.Context, userID int64) ([]Account, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT a.id, a.account_type_id, a.balance, a.is_system, a.created_at
		 FROM accounts a
		 INNER JOIN users_accounts ua ON ua.account_id = a.id
		 WHERE ua.user_id = $1
		 ORDER BY a.id`,
		userID)

	if err != nil {
		return nil, fmt.Errorf("select accounts: %w", err)
	}

	defer rows.Close()

	accounts, err := pgx.CollectRows(rows, pgx.RowToStructByName[Account])
	if err != nil {
		return nil, fmt.Errorf("select accounts: %w", err)
	}

	return accounts, nil
}

func (r *Repository) GetTransfers(ctx context.Context, userID int64) (*TransfersList, error) {
	rowsSent, err := r.pool.Query(ctx,
		`SELECT tr.id, tr.amount, tr.account_from_id, tr.account_to_id, tr.created_at
		 FROM transfers tr
		 INNER JOIN users_accounts ua ON ua.account_id = tr.account_from_id
		 WHERE ua.user_id = $1
		 ORDER BY tr.id DESC
		 `,
		userID)

	if err != nil {
		return nil, fmt.Errorf("select transfers: %w", err)
	}

	transfersSent, err := pgx.CollectRows(rowsSent, pgx.RowToStructByName[Transfer])
	if err != nil {
		return nil, fmt.Errorf("select transfers: %w", err)
	}

	rowsSent.Close()

	rowsIncoming, err := r.pool.Query(ctx,
		`SELECT tr.id, tr.amount, tr.account_from_id, tr.account_to_id, tr.created_at,
				a.is_system
		 FROM transfers tr
		 INNER JOIN users_accounts ua ON ua.account_id = tr.account_to_id
		 INNER JOIN accounts a ON a.id = tr.account_from_id
		 WHERE ua.user_id = $1
		 ORDER BY tr.id DESC`,
		userID)

	if err != nil {
		return nil, fmt.Errorf("select transfers: %w", err)
	}

	defer rowsIncoming.Close()

	transfersGet, err := pgx.CollectRows(rowsIncoming, pgx.RowToStructByName[TransferIncoming])
	if err != nil {
		return nil, fmt.Errorf("select transfers: %w", err)
	}

	return &TransfersList{
		TransfersSent:     transfersSent,
		TransfersIncoming: transfersGet,
	}, nil
}
