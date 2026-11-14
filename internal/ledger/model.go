package ledger

import "time"

type NewTransfer struct {
	UserID         int64
	AccountToID    int64
	AccountFromID  int64
	Amount         int64
	IsSystem       bool
	IdempotencyKey string
}

type NewDeposit struct {
	UserID         int64
	AccountToID    int64
	Amount         int64
	IdempotencyKey string
}

type Account struct {
	ID            int64     `db:"id" json:"id"`
	AccountTypeID int64     `db:"account_type_id" json:"account_type_id"`
	Balance       int64     `db:"balance" json:"balance"`
	IsSystem      bool      `db:"is_system" json:"is_system"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

type Transfer struct {
	ID            int64     `db:"id" json:"id"`
	Amount        int64     `db:"amount" json:"amount"`
	AccountFromID int64     `db:"account_from_id" json:"account_from_id"`
	AccountToID   int64     `db:"account_to_id" json:"account_to_id"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

type TransferIncoming struct {
	Transfer
	IsSystem bool `db:"is_system" json:"is_system"`
}

type TransfersList struct {
	TransfersSent     []Transfer         `json:"transfers_sent"`
	TransfersIncoming []TransferIncoming `json:"transfers_incoming"`
}
