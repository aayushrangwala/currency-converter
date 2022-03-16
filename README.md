# Currency Converter

This repository holds the application to convert the currency of all supported countries in the expected format.

This application provides the simple way to convert the currency by providing the input as country code, amount and the target country code in which we want to convert.

##### Note:
The exchange rates being used for the conversion can be referred from any provider. The provider can be specified in the conversion request.

By default, `CurrencyLayer` is the exchange rates provider which is being used for getting the supported currencies and their exchange rates.  

## API

[This](./api/proto/v1alpha1/currencyconverter/currency_converter_server.proto) is the `proto` file defining the API

There are primarily 3 APIs which this application supports

- ```ListExchangeRates``` is the API to list the exchange rates from an exchange rates provider and with a specific base currency. Default is `USD`.
- ```Convert``` is the main API which provides the conversion of the currency to the expected format.
- ```BatchConvert``` is the API which can be used to convert the currencies in Batch

These APIs are exposed in `REST` and `gRPC` format.

Below are the Rest APIs for the conversion
- `HTTP1.1 GET https://domain:port/v1alpha1/currency/convert`

Used for the conversion of the single currency from input country code and amount to the target country code


Example: `https://domain:port/v1alpha1/currency/convert?code='USD'&value='100'&to='EUR'&exchange_provider='yahoo'`

This will return the response as follows:

```json
{
  "converted": {
    "code": "EUR",
    "value": "80"
  },
  "from": {
    "code": "USD",
    "value": "100"
  },
  "exchange_rate": "0.8",
  "conversion_datetime": "xxxx",
  "exchange_rate_datetime": "xxxx"
}
```

- `HTTP1.1 POST https://domain:port/v1alpha1/batch/currency/convert`

Used for the batch conversion of currency from input country code and value to the target country code.

Example: `https://domain:port/v1alpha1/batch/currency/convert`

`Request Body`
```json
 [
   {
     "from": {
       "code": "USD",
       "value": "100"
     },
     "to": "EUR"
   },
  {
    "from": {
      "code": "EUR",
      "value": "99"
    },
    "to": "INR",
    "exchange_provider": "coingecko"
  }
 ]
```

`Response`

```json
[
  {
    "converted": {
      "code": "EUR",
      "value": "80"
    },
    "from": {
      "code": "USD",
      "value": "100"
    },
    "exchange_rate": "0.8",
    "conversion_datetime": "xxxx",
    "exchange_rate_datetime": "xxxx"
  },
  {
    "converted": {
      "code": "INR",
      "value": "890"
    },
    "from": {
      "code": "EUR",
      "value": "99"
    },
    "exchange_rate": "9",
    "conversion_datetime": "xxxx",
    "exchange_rate_datetime": "xxxx"
  }
]
```

- `HTTP1.1 GET https://domain:port/v1alpha1/currency/rates?exchange_provider='fixer'`

Returns the exchange rates for all the supported countries, with Base as `USD` currency code.

`Response`
```json
{
  "currencies": [
    {
      "code": "USD",
      "value": "1"
    },
    {
      "code": "EUR",
      "value": "1.2"
    },
    {
      "code": "INR",
      "value": "0.7"
    }
    ...
    ...
    ...
  ]
}
```

#### Exchange rates provider

We default the `CurrencyLayer` as the default exchange rates provider for our application. This can be changed in the conversion requests.

We support 6 exchange rates provider as of now

- [CoinGecko](https://www.coingecko.com/en/api/documentation)
- [CurrencyLayer](https://currencylayer.com/documentation) `Default`
- [Fixer](https://fixer.io/documentation)
- [Google](https://cloud.google.com/skus/exchange-rates)
- [Yahoo](https://finance.yahoo.com/currencies)

Any of the above 6 provider can be specified in the conversion/batch conversion/live-rates requests.

###### Note:
The provider name has to be in all `small case`

## Operations
The main flow to operate for the currency conversion

### Cache
This application uses an in-memory cache to store the exchange rates for the currency with base USD

The format of behaviour of the cache is `Cache Miss`. Initially when the server is started and cache is initialzed,
the cache is loaded with exchange rates using currencyLayer provider.

Each entry in the in-memory cache will have a `TTL` of 5 minutes.

On every conversion request if the exchange rate for conversion is not found, we fetch the rates from `any` of the provider which returned successfully first and update the cache as well.

We also store the supported currencies by the exchange rates provider in the cache.

### Flow

User will call any of the above APIs to convert, batch convert or get the live exchange rates.

API handler will get the exchange rates from in-memory cache. If there is a cache miss, we will fetch the live rates, update the cache and return the converted response.

The API responses also has the `conversion_datetime` and `exchange_rates_datetime` which gives an idea on when the currency is converted and when the rates were refreshed.

The API handler will use `pkg/converter` to perform the actual conversion logic.

## Background Jobs

Since the currency exchange rates gets updated throughout the 24 hours clock for different countries, we need our application to be updated with the live exchange rates so the conversion is returned with the minimal error offset.

Also, each entry in the cache, there is a default expiration time of `5 minutes` after which that excahnge rate is considered as cache miss.

We need an ability to cleanup those expired entries so that we dont clog the memory resources of the service. Also we store supported currencies by the provider in the cache, which has an expiration of 2 weeks.

To achieve the above cases, we introduced 2 background jobs:

- `Refresher` Whcih runs every 5 minutes and refreshes the exchange rates ion memory
- `Cleaner` runs every `5 minutes` to clean the entries which are expired. Also cleans the live supported currencies list whcih has expiry of 2 weeks.

###### Note:

Supported currencies are not refreshed as the part of background job but it refreshes as the part of cache miss.


# Development

### High Level Code Structure

- [API Proto](./api/proto/v1alpha1/currencyconverter/currency_converter_server.proto)

Proto file containing the specs for all the apis

- `internal`

This folder holds the utilities and modules which are private and not to be shared with the importers of this repository

- [inmemory](./internal/cache/inmemory)

The inmemory cache implements this [interface](./internal/cache/interface.go) interface. This interface can be implemented by other caches as well in future. Example: Redis as the store.

This folder holds the cache where the exchange rates and supported currencies will be stored and refreshed and cleaned.

- [exchange](./internal/exchange)

The exchange rates providers implements this [interface](./internal/exchange/interface.go) and can be implemented by any of the provider. 

This will hold the packages which implements the interface for a provider. Whenever we want to add a support of a new provider, we just have to add a new package in this path and implement the interface.

- [factory](./internal/factory)

The factory package is the provider for objects using factory design pattern

- [backgroundjobs](./pkg/backgroundjobs)

Holds the list of background jobs started in `main.go` which are refresher and cleaner.

- [server](./pkg/server)

Holds the API handler for the APIs defined in the proto

- [converter](./pkg/converter)

Will have the actual logic to convert the currency using the exchange rates from the cache

### Linter

We are using [`golangci-lint`](https://github.com/golangci/golangci-lint) as the linter for the code except test files

Run `make lint` to run the linter on the service

### Containerisation

We are using this [Dockerfile](./Dockerfile) for creating the images of the service

### Dev automation

We are using [`Makefile`](./Makefile) which has multiple targets to automate development and CI CD

### Testing

To run the unit test we have target in Makefile `make test` 