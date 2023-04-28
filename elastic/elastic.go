package elastic

import (
	"context"
	"fmt"

	"github.com/karatekaneen/crossrefindexer"
)

// Indexer ...
type Indexer struct {
	BatchSize int
}

// Index ...
func (i *Indexer) Index(
	ctx context.Context,
	data chan crossrefindexer.SimplifiedPublication,
) error {
	dataToIndex := make([]crossrefindexer.SimplifiedPublication, 0, i.BatchSize)
	for {
		// we receive a publication
		spub, open := <-data
		if !open {
			// If the channel is closed  - We "commit" the publications already in the slice before returning
			if err := indexPublications(ctx, dataToIndex); err != nil {
				return fmt.Errorf("Indexing failed: %w", err)
			}
			break
		}

		// we check if we have room in the slice
		if len(dataToIndex) < cap(dataToIndex) {
			// - If we do, then append
			// - fill up the slice if it's not full
			dataToIndex = append(dataToIndex, spub)
		}

		if len(dataToIndex) == cap(dataToIndex) {
			// - If we don't "commit" the slice to elasticsearch
			if err := indexPublications(ctx, dataToIndex); err != nil {
				return fmt.Errorf("Indexing failed: %w", err)
			}
			dataToIndex = dataToIndex[:0]
		}

	}

	return nil
}

func indexPublications(ctx context.Context, pubs []crossrefindexer.SimplifiedPublication) error {
	return nil // This is a placeholder
}
