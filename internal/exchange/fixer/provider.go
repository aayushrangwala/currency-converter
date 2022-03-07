package currencylayer

import "currency-converter/internal/exchange"

var _ exchange.Provider = (*provider)(nil)

type provider struct {
	api string
}

func New() exchange.Provider {
	return nil
}

func (p *provider) LiveRates() (map[string]string, error) {
	return nil, nil
}

func (p *provider) Currencies() ([]string, error) {
	return nil, nil
}
