package domain

import "errors"

var (
	ErrSybmolIsEmpty = errors.New("symbol is empty")
)

type Currency struct {
	Symbol string
}

func NewCurrency(symbol string) (*Currency, error) {
	if symbol == "" {
		return nil, ErrSybmolIsEmpty
	}

	return &Currency{Symbol: symbol}, nil
}
