package cache

import (
	//nolint:gosec
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	"currency-converter/internal/exchange"
)

const DefaultExpiration = 2 * time.Minute

// Store holds the currency exchange rates for an exchange provider with the currency code.
type Store interface {
	// AvailableCurrencies returns all the available currencies from the cache for an exchange provider. (only on Cache Miss or lazy populating)
	AvailableCurrencies(exchangeProvider exchange.ProviderType) ([]string, error)

	// SetAvailableCurrencies sets the available currencies to the cache for an exchange provider.
	SetAvailableCurrencies(exchangeProvider exchange.ProviderType, currencyCodes []string) error

	// GetExchangeRate returns the exchange rate for the passed currency code.
	// returns NotFound error if key not found.
	GetExchangeRate(currencyCode string, exchangeProvider exchange.ProviderType) (float32, error)

	// SetExchangeRate sets the exchange rate for a key from exchange provider and the currency code.
	// By default, each rate will have an expiration of 2 minutes.
	SetExchangeRate(currencyCode string, exchangeProvider exchange.ProviderType, rate float32, expiration time.Duration) error

	// RefreshExchangeRates fetches the latest exchange rates from all the supported exchange rates providers.
	// Will be used to refresh rates at:
	// 1. "Cache Miss" in a request for that provider.
	// 2. Every 5 minute refresh.
	RefreshExchangeRates(providers []exchange.ProviderType) error

	// CleanupAllExpired will cleanup all the entries which are expired.
	// One of its usage is going to be in a background job running every 5 minute.
	CleanupAllExpired()
}

// GetKey is the constructor for cache key using exchange provider and the currency code.
// It will be used to set and get the rates.
func GetKey(currencyCode string, exchangeProvider exchange.ProviderType) string {
	jsonBytes, _ := json.Marshal(struct {
		Code     string
		Provider string
	}{
		Code:     currencyCode,
		Provider: string(exchangeProvider),
	})
	//nolint:gosec
	md5Sum := md5.Sum(jsonBytes)
	return fmt.Sprintf("%x", md5Sum[:])
}
