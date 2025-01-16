package domain

import (
	"errors"
	"time"
)

var (
	ErrSybmolIsEmpty = errors.New("symbol is empty")
)

type Currency struct {
	Symbol string `db:"symbol"`
}

func NewCurrency(symbol string) (*Currency, error) {
	if symbol == "" {
		return nil, ErrSybmolIsEmpty
	}

	return &Currency{Symbol: symbol}, nil
}

type Price struct {
	Currency  Currency
	Price     float64
	Timestamp time.Time
}
