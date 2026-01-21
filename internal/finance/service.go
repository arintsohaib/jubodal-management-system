package finance

import (
	"context"
	"fmt"

	"github.com/bjdms/api/internal/models"
	"github.com/google/uuid"
)

// Service handles business logic for BJDMS finance
type Service struct {
	repo *Repository
}

// NewService creates a new finance service
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// RecordTransaction handles the creation of income/expense entries
func (s *Service) RecordTransaction(ctx context.Context, t *models.FinanceTransaction) error {
	// 1. Basic validation
	if t.Amount <= 0 {
		return fmt.Errorf("transaction amount must be greater than zero")
	}

	// 2. Expense check: Cannot spend more than current balance
	if t.Type == models.TransactionTypeExpense {
		balance, err := s.repo.GetBalance(ctx, t.JurisdictionID)
		if err != nil {
			return err
		}
		if balance.CurrentBalance < t.Amount {
			return fmt.Errorf("insufficient balance: current balance is Tk %.2f", balance.CurrentBalance)
		}
	}

	// 3. Persist (Immutability enforced by DB triggers)
	return s.repo.CreateTransaction(ctx, t)
}

// GetJurisdictionStatement returns a financial summary and recent transactions
func (s *Service) GetJurisdictionStatement(ctx context.Context, jurisdictionID uuid.UUID, page, pageSize int) (*models.FinanceBalance, []*models.FinanceTransaction, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	balance, err := s.repo.GetBalance(ctx, jurisdictionID)
	if err != nil {
		return nil, nil, err
	}

	transactions, err := s.repo.ListTransactions(ctx, jurisdictionID, pageSize, offset)
	if err != nil {
		return nil, nil, err
	}

	return balance, transactions, nil
}

// ListCategories returns available categories for recording
func (s *Service) ListCategories(ctx context.Context, transType string) ([]*models.FinanceCategory, error) {
	return s.repo.ListCategories(ctx, transType)
}
