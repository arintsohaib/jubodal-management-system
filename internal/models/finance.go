package models

import (
	"time"

	"github.com/google/uuid"
)

// Transaction types
const (
	TransactionTypeIncome   = "income"
	TransactionTypeExpense  = "expense"
	TransactionTypeTransfer = "transfer"
)

// FinanceCategory represents a category for a financial transaction
type FinanceCategory struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	NameBn    string    `json:"name_bn" db:"name_bn"`
	Type      string    `json:"type" db:"type"`
	IsSystem  bool      `json:"is_system" db:"is_system"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// FinanceTransaction represents a single financial entry in the ledger
type FinanceTransaction struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	JurisdictionID  uuid.UUID  `json:"jurisdiction_id" db:"jurisdiction_id"`
	UserID          uuid.UUID  `json:"user_id" db:"user_id"`
	CategoryID      uuid.UUID  `json:"category_id" db:"category_id"`
	Type            string     `json:"type" db:"type"`
	Amount          float64    `json:"amount" db:"amount"`
	Description     string     `json:"description" db:"description"`
	ReferenceNo     string     `json:"reference_no" db:"reference_no"`
	TransactionDate time.Time  `json:"transaction_date" db:"transaction_date"`
	EvidencePath    string     `json:"evidence_path" db:"evidence_path"`
	Metadata        any        `json:"metadata" db:"metadata"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`

	// Joined fields
	CategoryName   string `json:"category_name,omitempty" db:"category_name"`
	CategoryNameBn string `json:"category_name_bn,omitempty" db:"category_name_bn"`
	UserName       string `json:"user_name,omitempty" db:"user_name"`
}

// FinanceBalance represents the cached financial standing of a jurisdiction
type FinanceBalance struct {
	JurisdictionID uuid.UUID `json:"jurisdiction_id" db:"jurisdiction_id"`
	TotalIncome    float64   `json:"total_income" db:"total_income"`
	TotalExpense   float64   `json:"total_expense" db:"total_expense"`
	CurrentBalance float64   `json:"current_balance" db:"current_balance"`
	LastUpdatedAt  time.Time `json:"last_updated_at" db:"last_updated_at"`

	// Joined fields
	JurisdictionName string `json:"jurisdiction_name,omitempty" db:"jurisdiction_name"`
}
