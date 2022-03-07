package exchange

// Provider represents different exchange providers.
type Provider interface {
	// LiveRates fetch the live exchange rates for all the supported currencies.
	LiveRates() (map[string]float32, error)

	// Currencies lists all the available/supported currencies by the provider.
	Currencies() ([]string, error)
}

// ProviderType represents the type fo the exchange rates provider supported.
type ProviderType string

const (
	CurrencyLayer = "currencylayer"

	CoinGecko = "coingecko"

	Google = "google"

	Fixer = "fixer"

	OpenExchangeRates = "openexchangerates"

	Yahoo = "yahoo"
)

// GetSupportedProviders returns the list of supported exchange rate providers.
func GetSupportedProviders() []ProviderType {
	return []ProviderType{
		CurrencyLayer,
		CoinGecko,
		Google,
		Fixer,
		OpenExchangeRates,
		Yahoo,
	}
}
