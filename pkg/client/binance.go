package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sixojke/crypto-service/internal/domain"
	"github.com/sixojke/crypto-service/pkg/logger"
)

const getPriceAPI = "https://api.binance.com/api/v3/ticker/price?symbol="

type Price struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price,string"`
}

// GetPrice retrieves the current price of a cryptocurrency from Binance.
func GetPrice(ctx context.Context, symbol string) (*domain.Price, error) {
	url := getPriceAPI + strings.ToUpper(symbol)
	logger.Error(url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var price Price
	if err := json.NewDecoder(resp.Body).Decode(&price); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &domain.Price{
		Currency:  domain.Currency{Symbol: price.Symbol},
		Price:     price.Price,
		Timestamp: time.Now(),
	}, nil
}
