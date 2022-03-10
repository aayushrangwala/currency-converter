package backgroundjobs

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"currency-converter/internal/cache"
	"currency-converter/internal/exchange"
)

func RunExchangeRatesRefresher(ctx context.Context, store cache.Store, interval time.Duration) {
	// refreshes the exchange rates for the default or first successful provider
	var cancelFunc context.CancelFunc
	ctx, cancelFunc = signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer cancelFunc()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := store.RefreshExchangeRates(exchange.GetSupportedProviders()); err != nil {
				logrus.WithError(err).Error("failed to refresh exchange rates")
			}

		case <-ctx.Done():
			cancelFunc()
			return
		}
	}
}
