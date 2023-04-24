package elastic

import (
	"context"
	"fmt"

	"github.com/karatekaneen/crossrefindexer"
)

type Indexer struct {
	dataToIndex []crossrefindexer.SimplifiedPublication
}

func (i *Indexer) Index(
	ctx context.Context,
	data chan crossrefindexer.SimplifiedPublication,
) error {
	// TODO:
	// - fill up the slice if it's not full

	// we receive a publication
	// we check if we have room in the slice
	// - If we do, then append
	// - If we don't "commit" the slice to elasticsearch
	// If the channel is closed  - We "commit" the publications already in the slice before returning
	return fmt.Errorf("not implemented")
}

func indexPublications(ctx context.Context, pubs []crossrefindexer.SimplifiedPublication) error {
	return nil // This is a placeholder
}
