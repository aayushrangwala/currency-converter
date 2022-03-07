package currencylayer

import (
	"currency-converter/internal/exchange"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ exchange.Provider = (*provider)(nil)

type provider struct {
	api string
}

func New() exchange.Provider {
	return nil
}

func (p *provider) LiveRates() (map[string]float32, error) {
	return nil, status.Error(codes.Unimplemented, "function not implemented for the provider")
}

func (p *provider) Currencies() ([]string, error) {
	return nil, status.Error(codes.Unimplemented, "function not implemented for the provider")
}
