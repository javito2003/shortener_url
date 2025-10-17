package clicks_worker

import (
	"context"
	"log"
	"time"
)

// Service orquesta la lógica del worker.
type Service struct {
	reader  ClickCacheReader
	updater LinkBulkUpdater
}

func NewService(r ClickCacheReader, u LinkBulkUpdater) *Service {
	return &Service{reader: r, updater: u}
}

func (s *Service) Run(ctx context.Context, interval time.Duration) {
	log.Println("Clicks worker started...")
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Worker tick: processing clicks...")
			if err := s.processClicks(ctx); err != nil {
				log.Printf("ERROR: Failed to process clicks: %v", err)
			}
		case <-ctx.Done():
			log.Println("Clicks worker shutting down...")
			return
		}
	}
}

// processClicks es la lógica de un solo ciclo.
func (s *Service) processClicks(ctx context.Context) error {
	clickCounts, err := s.reader.FetchAndClear(ctx)
	if err != nil {
		return err
	}

	if len(clickCounts) == 0 {
		log.Println("No new clicks to process.")
		return nil
	}

	log.Printf("Processing %d links with new clicks.", len(clickCounts))

	if err := s.updater.IncrementClickCounts(ctx, clickCounts); err != nil {
		return err
	}

	log.Println("Successfully updated click counts in the database.")
	return nil
}
