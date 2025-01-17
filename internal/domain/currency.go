package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/sixojke/crypto-service/pkg/binance"
)

var (
	ErrSybmolIsEmpty       = errors.New("symbol is empty")
	ErrSymbolDoesNotExists = errors.New("symbol does not exists")
)

type Currency struct {
	Symbol string `db:"symbol"`
}

func NewCurrency(symbol string) (*Currency, error) {
	if symbol == "" {
		return nil, ErrSybmolIsEmpty
	}

	upperSymbol := strings.ToUpper(symbol)

	ok, err := binance.CheckSymbol(upperSymbol)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, ErrSymbolDoesNotExists
	}

	return &Currency{Symbol: upperSymbol}, nil
}

type Price struct {
	Currency  Currency
	Price     float64
	Timestamp time.Time
}
