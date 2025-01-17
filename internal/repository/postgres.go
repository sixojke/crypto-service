package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/sixojke/crypto-service/internal/domain"
	"github.com/sixojke/crypto-service/pkg/logger"
)

type Currency interface {
	AddTrackedCurrency(currency *domain.Currency) error
	SavePrice(price *domain.Price) error
	GetTrackedCurrencies() ([]domain.Currency, error)
	RemoveFromTracking(symbol string) error
	GetPriceByTimestamp(symbol string, timestamp time.Time) (*domain.Price, error)
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

		logger.Error(err.Error())
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

	if _, err := r.DB.Exec(query, price.Currency, price.Price, price.Timestamp); err != nil {
		logger.Error(err.Error())
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
		logger.Error(err.Error())
		return []domain.Currency{}, fmt.Errorf("failed to get tracked currencies: %v", err)
	}

	return currencies, nil
}

func (r *CurrencyPostgres) RemoveFromTracking(symbol string) error {
	query := `
	DELETE FROM 
		tracked_currencies
	WHERE
		symbol = $1
	`

	if _, err := r.DB.Exec(query, symbol); err != nil {
		logger.Error(err.Error())
		return fmt.Errorf("failed to remove currency from tracking: symbol=%v", symbol)
	}

	return nil
}

func (r *CurrencyPostgres) GetPriceByTimestamp(symbol string, timestamp time.Time) (*domain.Price, error) {
	query := `
		SELECT 
			symbol, price, timestamp 
		FROM 
			currency_prices
		WHERE 
			symbol = $1 
			AND timestamp <= $2
		ORDER BY 
			timestamp DESC
		LIMIT 1;
	`

	var price domain.Price
	if err := r.DB.Get(&price, query, symbol, timestamp); err != nil {
		// If no values ​​were found before the requested timestamp, look after
		if err == sql.ErrNoRows {
			query = `
				SELECT 
					symbol, price, timestamp
				FROM 
					currency_prices
				WHERE symbol = $1 AND timestamp > $2
				ORDER BY 
					timestamp ASC
				LIMIT 1;
			`

			if err := r.DB.Get(&price, query, symbol, timestamp); err != nil {
				// Data does not exist for this currency
				if err == sql.ErrNoRows {
					return nil, domain.ErrNoDataOnThisCurrency
				}

				logger.Error(err.Error())
				return nil, fmt.Errorf("failed to get price: %v", err)
			} else {
				return &price, nil
			}

		} else {
			logger.Error(err.Error())
			return nil, fmt.Errorf("failed to get price: %v", err)
		}
	}

	return &price, nil
}
