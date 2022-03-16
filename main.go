package main

import (
	"context"
	"log"
	"time"

	"golang.org/x/sync/errgroup"

	"currency-converter/internal/cache/inmemory"
	"currency-converter/pkg/backgroundjobs"
)

func main() {
	// start background jobs
	// start the servers

	g, ctx := errgroup.WithContext(context.Background())

	store := inmemory.NewStore()

	g.Go(func() error {
		backgroundjobs.RunCacheCleaner(ctx, store, 5*time.Minute)
		return nil
	})

	g.Go(func() error {
		backgroundjobs.RunExchangeRatesRefresher(ctx, store, 5*time.Minute)
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
