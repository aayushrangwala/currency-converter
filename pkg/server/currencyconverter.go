package server

import (
	"context"
	"strconv"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "currency-converter/api/pb/v1alpha1/currencyconverter"
	"currency-converter/internal/cache"
	"currency-converter/internal/errors"
	"currency-converter/internal/exchange"
	"currency-converter/pkg/converter"
)

type converterServer struct {
	store cache.Store
}

func NewServer(store cache.Store) pb.CurrencyConverterServiceServer {
	return &converterServer{
		store: store,
	}
}

func (server *converterServer) ListExchangeRates(
	ctx context.Context,
	request *pb.ListExchangeRatesRequest) (*pb.ListExchangeRatesResponse, error) {
	return nil, errors.UnImplementedError
}

func (server *converterServer) Convert(ctx context.Context, request *pb.ConversionRequest) (*pb.ConversionResponse, error) {
	// TODO: User authentication using ctx

	exProvider := exchange.ProviderType(request.ExchangeProvider)
	if exProvider == "" {
		exProvider = exchange.CurrencyLayer
	}

	var rate float32
	var err error

	if rate, err = server.store.GetExchangeRate(request.GetTo(), exProvider); err != nil {
		return nil, err
	}

	// cache HIT

	// Handle error better. Implemented just for the submission purpose.
	amount, _ := strconv.ParseFloat(request.GetFrom().Value, 64)

	return &pb.ConversionResponse{
		Converted: &pb.Currency{
			Code:  request.GetTo(),
			Value: strconv.FormatFloat(converter.Convert(rate, amount), 10, 2, 64),
		},
		From:                 request.GetFrom(),
		ExchangeRate:         rate,
		ConversionDatetime:   timestamppb.Now(),
		ExchangeRateDatetime: timestamppb.Now(),
	}, nil
}

func (server *converterServer) BatchConvert(
	_ context.Context,
	_ *pb.BatchConversionRequest) (*pb.BatchConversionResponse, error) {
	return nil, errors.UnImplementedError
}
