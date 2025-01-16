package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sixojke/crypto-service/internal/domain"
)

type Currency interface {
	Add(currency *domain.Currency) error
}

type CurrencyPostgres struct {
	DB *sqlx.DB
}

// NewCurrencyPostgres creates a new currency repository.
func NewCurrencyPostgres(db *sqlx.DB) Currency {
	return &CurrencyPostgres{DB: db}
}

func (r *CurrencyPostgres) Add(currency *domain.Currency) error {
	query := `INSERT INTO tracked_currencies (symbol) VALUES ($1)`

	if _, err := r.DB.Exec(query, currency.Symbol); err != nil {
		return fmt.Errorf("failed to add currency: %v", err)
	}

	return nil
}
