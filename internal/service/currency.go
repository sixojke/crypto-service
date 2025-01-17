package service

import (
	"context"
	"slices"
	"strings"
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

// RemoveFromTracking removes a currency from the tracking list.
func (s *CurrencyService) RemoveFromTracking(symbol string) error {
	upperSymbol := strings.ToUpper(symbol)

	ok, err := s.repo.RemoveFromTracking(strings.ToUpper(upperSymbol))
	if err != nil {
		return err
	}

	logger.Errorf("%v", trackedCurrencies)
	if ok {
		for i, cur := range trackedCurrencies {
			if upperSymbol == cur.Symbol {
				trackedCurrencies = slices.Delete(trackedCurrencies, i, i+1)
				break
			}
		}
	}
	logger.Errorf("%v", trackedCurrencies)

	return nil
}

// GetPriceByTimestamp retrieves the price for the given symbol at the specified Unix timestamp.
func (s *CurrencyService) GetPriceByTimestamp(symbol string, unixTimestamp int64) (*domain.Price, error) {
	return s.repo.GetPriceByTimestamp(strings.ToUpper(symbol), time.Unix(unixTimestamp, 0))
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
					Currency:  price.Currency,
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
