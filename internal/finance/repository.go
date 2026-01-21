package finance

import (
	"context"
	"time"

	"github.com/bjdms/api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database operations for finance
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new finance repository
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

// CreateCategory adds a new transaction category
func (r *Repository) CreateCategory(ctx context.Context, c *models.FinanceCategory) error {
	query := `INSERT INTO finance_categories (name, name_bn, type, is_system) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, c.Name, c.NameBn, c.Type, c.IsSystem).Scan(&c.ID, &c.CreatedAt)
}

// ListCategories returns all transaction categories
func (r *Repository) ListCategories(ctx context.Context, transType string) ([]*models.FinanceCategory, error) {
	query := `SELECT id, name, name_bn, type, is_system, created_at FROM finance_categories`
	var args []interface{}
	if transType != "" {
		query += " WHERE type = $1"
		args = append(args, transType)
	}
	query += " ORDER BY name ASC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.FinanceCategory
	for rows.Next() {
		var c models.FinanceCategory
		if err := rows.Scan(&c.ID, &c.Name, &c.NameBn, &c.Type, &c.IsSystem, &c.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, &c)
	}
	return list, nil
}

// CreateTransaction inserts a new financial ledger entry
func (r *Repository) CreateTransaction(ctx context.Context, t *models.FinanceTransaction) error {
	// Transactions are immutable, balance update is handled by DB trigger
	query := `
		INSERT INTO finance_transactions (
			jurisdiction_id, user_id, category_id, type, amount, 
			description, reference_no, transaction_date, evidence_path, metadata
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at
	`
	if t.TransactionDate.IsZero() {
		t.TransactionDate = time.Now()
	}
	return r.db.QueryRow(ctx, query,
		t.JurisdictionID, t.UserID, t.CategoryID, t.Type, t.Amount,
		t.Description, t.ReferenceNo, t.TransactionDate, t.EvidencePath, t.Metadata,
	).Scan(&t.ID, &t.CreatedAt)
}

// GetBalance retrieves the current financial standing of a jurisdiction
func (r *Repository) GetBalance(ctx context.Context, jurisdictionID uuid.UUID) (*models.FinanceBalance, error) {
	query := `
		SELECT fb.jurisdiction_id, fb.total_income, fb.total_expense, fb.current_balance, fb.last_updated_at,
		       j.name as jurisdiction_name
		FROM finance_balances fb
		JOIN jurisdictions j ON fb.jurisdiction_id = j.id
		WHERE fb.jurisdiction_id = $1
	`
	var b models.FinanceBalance
	err := r.db.QueryRow(ctx, query, jurisdictionID).Scan(
		&b.JurisdictionID, &b.TotalIncome, &b.TotalExpense, &b.CurrentBalance, &b.LastUpdatedAt,
		&b.JurisdictionName,
	)
	if err != nil {
		// New jurisdictions might not have a balance record yet
		return &models.FinanceBalance{JurisdictionID: jurisdictionID}, nil
	}
	return &b, nil
}

// ListTransactions returns financial activity for a jurisdiction
func (r *Repository) ListTransactions(ctx context.Context, jurisdictionID uuid.UUID, limit, offset int) ([]*models.FinanceTransaction, error) {
	query := `
		SELECT ft.id, ft.jurisdiction_id, ft.user_id, ft.category_id, ft.type, ft.amount, 
		       ft.description, ft.reference_no, ft.transaction_date, ft.evidence_path, ft.created_at,
		       fc.name as category_name, fc.name_bn as category_name_bn, u.full_name as user_name
		FROM finance_transactions ft
		JOIN finance_categories fc ON ft.category_id = fc.id
		JOIN users u ON ft.user_id = u.id
		WHERE ft.jurisdiction_id = $1
		ORDER BY ft.transaction_date DESC, ft.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, jurisdictionID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.FinanceTransaction
	for rows.Next() {
		var t models.FinanceTransaction
		err := rows.Scan(
			&t.ID, &t.JurisdictionID, &t.UserID, &t.CategoryID, &t.Type, &t.Amount,
			&t.Description, &t.ReferenceNo, &t.TransactionDate, &t.EvidencePath, &t.CreatedAt,
			&t.CategoryName, &t.CategoryNameBn, &t.UserName,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, &t)
	}
	return list, nil
}
