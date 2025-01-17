package repository

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/sixojke/crypto-service/internal/domain"
)

type Currency interface {
	AddTrackedCurrency(currency *domain.Currency) error
	SavePrice(price *domain.Price) error
	GetTrackedCurrencies() ([]domain.Currency, error)
}

type CurrencyPostgres struct {
	DB *sqlx.DB
}

// NewCurrencyPostgres creates a new currency repository.
func NewCurrencyPostgres(db *sqlx.DB) Currency {
	return &CurrencyPostgres{DB: db}
}

func (r *CurrencyPostgres) AddTrackedCurrency(currency *domain.Currency) error {
	query := `INSERT INTO tracked_currencies (symbol) VALUES ($1)`

	if _, err := r.DB.Exec(query, currency.Symbol); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil
		}

		return fmt.Errorf("failed to add currency: %v", err)
	}

	return nil
}

func (r *CurrencyPostgres) SavePrice(price *domain.Price) error {
	query := `
		INSERT INTO 
			currency_prices
		(
			symbol,
			price,
			timestamp
		) VALUES ($1, $2, $3)
	`

	if _, err := r.DB.Exec(query, price.Currency.Symbol, price.Price, price.Timestamp); err != nil {
		return fmt.Errorf("failed to save price: %v", err)
	}

	return nil
}

func (r *CurrencyPostgres) GetTrackedCurrencies() ([]domain.Currency, error) {
	query := `
		SELECT
			symbol
		FROM 
			tracked_currencies
	`

	currencies := make([]domain.Currency, 0)
	if err := r.DB.Select(&currencies, query); err != nil {
		return []domain.Currency{}, fmt.Errorf("failed to get tracked currencies: %v", err)
	}

	return currencies, nil
}
