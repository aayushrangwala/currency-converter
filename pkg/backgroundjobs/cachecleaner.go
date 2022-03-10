package backgroundjobs

import (
	"context"
	"os/signal"
	"syscall"
	"time"

	"currency-converter/internal/cache"
)

func RunCacheCleaner(ctx context.Context, store cache.Store, interval time.Duration) {
	var cancelFunc context.CancelFunc
	ctx, cancelFunc = signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer cancelFunc()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			store.CleanupAllExpired()

		case <-ctx.Done():
			cancelFunc()
			return
		}
	}
}
