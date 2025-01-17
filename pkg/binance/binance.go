package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const getPriceAPI = "https://api.binance.com/api/v3/ticker/price?symbol="

type Price struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price,string"`
}

// GetPrice retrieves the current price of a cryptocurrency from Binance.
func GetPrice(ctx context.Context, symbol string) (*Price, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, generateGetPriceUrl(symbol), nil)
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

	return &Price{
		Symbol: symbol,
		Price:  price.Price,
	}, nil
}

func CheckSymbol(symbol string) (bool, error) {
	resp, err := http.Get(generateGetPriceUrl(symbol))
	if err != nil {
		return false, fmt.Errorf("creating request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	return true, nil
}

func generateGetPriceUrl(symbol string) string {
	return getPriceAPI + symbol
}
