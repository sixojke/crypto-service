package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/sixojke/crypto-service/pkg/binance"
	"github.com/sixojke/crypto-service/pkg/logger"
)

var (
	ErrSybmolIsEmpty        = errors.New("symbol is empty")
	ErrSymbolDoesNotExists  = errors.New("symbol does not exists")
	ErrNoDataOnThisCurrency = errors.New("no data on this currency")
	ErrInvalidTimestamp     = errors.New("invalid timestamp")
)

type Currency struct {
	Symbol string `db:"symbol"`
}

func NewCurrency(symbol string) (*Currency, error) {
	if symbol == "" {
		return nil, ErrSybmolIsEmpty
	}

	upperSymbol := strings.ToUpper(symbol)

	ok, err := binance.CheckCurrency(upperSymbol)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	if !ok {
		return nil, ErrSymbolDoesNotExists
	}

	return &Currency{Symbol: upperSymbol}, nil
}

type Price struct {
	Currency  string    `db:"symbol"`
	Price     float64   `db:"price"`
	Timestamp time.Time `db:"timestamp"`
}
