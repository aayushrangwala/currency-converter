package factory

import (
	"currency-converter/internal/exchange"
	"currency-converter/internal/exchange/coingecko"
	"currency-converter/internal/exchange/currencylayer"
	"currency-converter/internal/exchange/fixer"
	"currency-converter/internal/exchange/google"
	"currency-converter/internal/exchange/openexchangerates"
	"currency-converter/internal/exchange/yahoo"
)

type exchangeRatesProviderFactory struct {
}

func NewExchangeRatesProviderFactory() *exchangeRatesProviderFactory {
	return &exchangeRatesProviderFactory{}
}

func (factory *exchangeRatesProviderFactory) BuildExchangeRatesProvider(providerType exchange.ProviderType) exchange.Provider {
	var provider exchange.Provider

	switch providerType {
	case exchange.CurrencyLayer:
		provider = currencylayer.New()
	case exchange.CoinGecko:
		provider = coingecko.New()
	case exchange.Google:
		provider = google.New()
	case exchange.Fixer:
		provider = fixer.New()
	case exchange.OpenExchangeRates:
		provider = openexchangerates.New()
	case exchange.Yahoo:
		provider = yahoo.New()
	}

	return provider
}
