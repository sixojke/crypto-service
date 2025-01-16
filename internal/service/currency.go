package service

import (
	"github.com/sixojke/crypto-service/internal/domain"
	"github.com/sixojke/crypto-service/internal/repository"
)

type CurrencyService struct {
	repo repository.Currency
}

func NewCurrencyService(repo repository.Currency) *CurrencyService {
	return &CurrencyService{
		repo: repo,
	}
}

func (s *CurrencyService) Add(symbol string) error {
	currency, err := domain.NewCurrency(symbol)
	if err != nil {
		return err
	}

	if err := s.repo.Add(currency); err != nil {
		return err
	}

	return nil
}
