package inmemory

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"currency-converter/internal/cache"
	apierrs "currency-converter/internal/errors"
	"currency-converter/internal/exchange"
	"currency-converter/internal/factory"
)

var _ cache.Store = (*inMemory)(nil)

// entry represents an entry of the inMemory cache
type entry struct {
	// data is the actual content to be stored in the cache.
	data interface{}

	// expiration is the calculated time based on the passed validity, till this entry is valid.
	expiration time.Time
}

// newStoreEntry returns an entry from the data and validity, which is to be stored in the cache.
func newStoreEntry(data interface{}, validity time.Duration) *entry {
	e := &entry{
		data: data,
	}

	if validity > 0 {
		e.expiration = time.Now().Add(validity)
	}

	return e
}

// IsExpired checks whether the entry stored is expired.
func (e *entry) IsExpired() bool {
	if e.expiration.IsZero() {
		// no expiration set
		return false
	}

	if time.Now().UnixNano() < e.expiration.UnixNano() {
		return false
	}

	return true
}

// inMemory is the cache store where the data will be in-memory.
type inMemory struct {
	items map[string]*entry
	mu    *sync.RWMutex
}

// NewStore is a constructor for inMemory cache store.
func NewStore() cache.Store {
	return &inMemory{
		items: map[string]*entry{},
		mu:    &sync.RWMutex{},
	}
}

// AvailableCurrencies returns the available currencies from cache for an exchange provider.
// Updates the cache for cache miss.
func (store *inMemory) AvailableCurrencies(exchangeProvider exchange.ProviderType) ([]string, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	if val, present := store.items[string(exchangeProvider)]; present {
		return val.data.([]string), nil
	}

	var err error
	var currencies []string

	// cache miss, update the cache and return the value
	provider := factory.NewExchangeRatesProviderFactory().BuildExchangeRatesProvider(exchangeProvider)
	if currencies, err = provider.Currencies(); err != nil {
		return []string{}, err
	}

	return currencies, store.SetAvailableCurrencies(exchangeProvider, currencies)
}

// SetAvailableCurrencies sets the list of available currencies from a provider to redis.
func (store *inMemory) SetAvailableCurrencies(exchangeProvider exchange.ProviderType, currencyCodes []string) error {
	if currencyCodes == nil || len(currencyCodes) == 0 {
		return apierrs.InvalidArgumentError
	}

	store.mu.Lock()
	defer store.mu.Unlock()

	store.items[string(exchangeProvider)] = newStoreEntry(currencyCodes, 2*7*24*time.Hour) // 2 weeks of validity

	return nil
}

// GetExchangeRate returns exchange rate for the passed currency code for the exchange provider.
func (store *inMemory) GetExchangeRate(currencyCode string, exchangeProvider exchange.ProviderType) (float32, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	key := cache.GetKey(currencyCode, exchangeProvider)

	if val, present := store.items[key]; present {
		return val.data.(float32), nil
	}

	// cache MISS: refresh rates in cache
	provider := factory.NewExchangeRatesProviderFactory().BuildExchangeRatesProvider(exchangeProvider)

	var err error
	var rates map[string]float32

	if rates, err = provider.LiveRates(); err != nil {
		return -1, apierrs.UpstreamExchangeRateServerError
	}

	if err = store.SetExchangeRate(currencyCode, exchangeProvider, rates[currencyCode], cache.DefaultExpiration); err != nil {
		return -1, apierrs.InternalCacheError
	}

	val, present := rates[currencyCode]
	if !present {
		return -1, apierrs.CacheKeyNotFoundError
	}

	return val, nil
}

func (store *inMemory) SetExchangeRate(
	currencyCode string,
	exchangeProvider exchange.ProviderType,
	rate float32,
	expiration time.Duration) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	key := cache.GetKey(currencyCode, exchangeProvider)

	store.items[key] = newStoreEntry(rate, expiration)

	return nil
}

func (store *inMemory) RefreshExchangeRates(providers []exchange.ProviderType) error {
	logger := logrus.New()

	g, _ := errgroup.WithContext(context.Background())

	for _, providerType := range providers {
		exchangeProvider := providerType
		g.Go(func() error {
			store.mu.Lock()
			defer store.mu.Unlock()

			provider := factory.NewExchangeRatesProviderFactory().BuildExchangeRatesProvider(exchangeProvider)

			var err error
			var rates map[string]float32

			if rates, err = provider.LiveRates(); err != nil {
				logger.WithError(err).Warnf("error while fetching live rates from the provider: [%s]", exchangeProvider)

				// returning with nil error as we want just one successful hit of the LiveRates from any of the passed provider.
				return nil
			}

			for code, rate := range rates {
				if tErr := store.SetExchangeRate(code, exchangeProvider, rate, cache.DefaultExpiration); tErr != nil {
					logger.WithError(err).
						Warnf("failed to set exchangerate for currency [%s] and provider [%s]", code, exchangeProvider)

					err = tErr
				}
			}

			if err != nil {
				// we send error only when we have the successful completion of live rates update.
				return nil
			}

			return status.Error(codes.OK, "error for sending the signal of first successful live rates updates")
		})
	}

	err := g.Wait()
	if err != nil && status.Code(err) == codes.OK {
		logger.Info("Successful live rates are updated by the one provider")

		return nil
	}

	logger.Error("Refresher workers stopped. No workers could update exchange rates from the given list of providers")
	return err
}

// CleanupAllExpired will delete all the expired entries.
func (store *inMemory) CleanupAllExpired() {
	store.mu.Lock()
	defer store.mu.Lock()

	for key, cacheEntry := range store.items {
		if !cacheEntry.IsExpired() {
			continue
		}

		delete(store.items, key)
	}
}
