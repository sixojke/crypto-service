package service

import (
	"context"
	"sync"
	"time"

	"github.com/sixojke/crypto-service/internal/config"
	"github.com/sixojke/crypto-service/internal/domain"
	"github.com/sixojke/crypto-service/internal/repository"
	"github.com/sixojke/crypto-service/pkg/binance"
	"github.com/sixojke/crypto-service/pkg/logger"
)

var trackedCurrencies = []domain.Currency{}

type CurrencyService struct {
	repo   repository.Currency
	config config.CurrencyService
}

func NewCurrencyService(repo repository.Currency, config config.CurrencyService) *CurrencyService {
	return &CurrencyService{
		repo:   repo,
		config: config,
	}
}

// AddToTracking adds a new currency symbol to the watch list
func (s *CurrencyService) AddToTracking(symbol string) error {
	currency, err := domain.NewCurrency(symbol)
	if err != nil {
		return err
	}

	if err := s.repo.AddTrackedCurrency(currency); err != nil {
		return err
	}

	trackedCurrencies = append(trackedCurrencies, domain.Currency{
		Symbol: currency.Symbol,
	})

	return nil
}

// LaunchCurrencyTracking starts the process of tracking currency prices.
func (s *CurrencyService) LaunchCurrencyTracking() {
	currencies, err := s.repo.GetTrackedCurrencies()
	if err != nil {
		logger.Fatal(err.Error())
	}

	trackedCurrencies = currencies

	interval := time.Duration(1000/s.config.UpdatesPerSercond) * time.Millisecond

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		concurrency := 10
		semaphore := make(chan struct{}, concurrency)

		var wg sync.WaitGroup
		for _, currency := range trackedCurrencies {
			wg.Add(1)
			go func(currency domain.Currency) {
				defer wg.Done()

				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				price, err := binance.GetPrice(ctx, currency.Symbol)
				if err != nil {
					if ctx.Err() != nil {
						logger.Warnf("GetPrice for %s timed out or canceled: %v", currency.Symbol, ctx.Err())
						return
					}
					logger.Errorf("Error getting price for %s: %v", currency.Symbol, err)
					return
				}

				if err := s.repo.SavePrice(&domain.Price{
					Currency: domain.Currency{
						Symbol: price.Symbol,
					},
					Price:     price.Price,
					Timestamp: time.Now(),
				}); err != nil {
					logger.Errorf("Error getting price for %s: %v", currency.Symbol, err)
					return
				}

			}(currency)
		}
		wg.Wait()
	}
}
