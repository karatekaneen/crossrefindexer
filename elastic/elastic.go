package elastic

import (
	"context"
	"fmt"
	"github.com/karatekaneen/crossrefindexer"
)

// Indexer ...
type Indexer struct {
	dataToIndex []crossrefindexer.SimplifiedPublication
}

// Index ...
func (i *Indexer) Index(ctx context.Context, data chan crossrefindexer.SimplifiedPublication) error {
	i.dataToIndex = make([]crossrefindexer.SimplifiedPublication, 0, 1000)
	for {
		// we receive a publication
		spub, open := <-data
		if !open {
			// If the channel is closed  - We "commit" the publications already in the slice before returning
			fmt.Println(i.dataToIndex)
			break
		}

		// we check if we have room in the slice
		if len(i.dataToIndex) < cap(i.dataToIndex) {
			// - If we do, then append
			// - fill up the slice if it's not full
			i.dataToIndex = append(i.dataToIndex, spub)
		}

		if len(i.dataToIndex) == cap(i.dataToIndex) {
			// - If we don't "commit" the slice to elasticsearch
			fmt.Println(i.dataToIndex)
			i.dataToIndex = i.dataToIndex[:0]
		}

	}

	return fmt.Errorf("not implemented")
}

func indexPublications(ctx context.Context, pubs []crossrefindexer.SimplifiedPublication) error {
	return nil // This is a placeholder
}
