package ledger

import (
	"context"
)

type LedgerRepository interface {
	CreateAccount(ctx context.Context, userID int64) (int64, error)
	Deposit(ctx context.Context, deposit NewDeposit) error
	Transfer(ctx context.Context, transfer NewTransfer) error
	GetAccounts(ctx context.Context, userID int64) ([]Account, error)
	GetTransfers(ctx context.Context, userID int64) (*TransfersList, error)
}

type Service struct {
	repo LedgerRepository
}

func NewService(repo LedgerRepository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateAccount(ctx context.Context, userID int64) (int64, error) {
	accountID, err := s.repo.CreateAccount(ctx, userID)

	if err != nil {
		return 0, err
	}

	return accountID, nil
}

func (s *Service) CreateDeposit(ctx context.Context, deposit NewDeposit) error {
	err := s.repo.Deposit(ctx, deposit)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) CreateTransfer(ctx context.Context, transfer NewTransfer) error {
	err := s.repo.Transfer(ctx, transfer)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetAccounts(ctx context.Context, userID int64) ([]Account, error) {
	accounts, err := s.repo.GetAccounts(ctx, userID)

	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *Service) GetTransfers(ctx context.Context, userID int64) (*TransfersList, error) {
	transfers, err := s.repo.GetTransfers(ctx, userID)

	if err != nil {
		return nil, err
	}

	return transfers, nil
}
