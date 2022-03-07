package exchange

import (
	"currency-converter/internal/exchange/coingecko"
	"currency-converter/internal/exchange/currencylayer"
	"currency-converter/internal/exchange/fixer"
	"currency-converter/internal/exchange/google"
	"currency-converter/internal/exchange/openexchangerates"
	"currency-converter/internal/exchange/yahoo"
)

type ProviderFactory struct {
}

func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{}
}

func (factory *ProviderFactory) BuildExchangeRatesProvider(providerType ProviderType) Provider {
	var provider Provider

	switch providerType {
	case CurrencyLayer:
		provider = currencylayer.New()
	case CoinGecko:
		provider = coingecko.New()
	case Google:
		provider = google.New()
	case Fixer:
		provider = fixer.New()
	case OpenExchangeRates:
		provider = openexchangerates.New()
	case Yahoo:
		provider = yahoo.New()
	}

	return provider
}
