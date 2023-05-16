package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	"github.com/karatekaneen/crossrefindexer"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// TODO: Add loading from env variables
type Config struct {
	BatchSize           int
	IndexName           string
	FlushBytes          int
	FlushInterval       int
	NumWorkers          int
	Password            string
	Username            string
	Addresses           []string
	CACert              []byte
	DisableRetry        bool
	MaxRetries          int
	CompressRequestBody bool
}

type Indexer struct {
	config Config
	client *elasticsearch.Client
	log    *zap.SugaredLogger
}

func New(config Config, log *zap.SugaredLogger) (*Indexer, error) {
	// Instantiate the exponential backoff thingy
	retryBackoff := backoff.NewExponentialBackOff()

	esClient, err := createElasticClient(config, retryBackoff)
	if err != nil {
		return nil, err
	}

	idx := &Indexer{
		config: config,
		client: esClient,
		log:    log,
	}

	return idx, nil
}

// IndexPublications ...
func (i *Indexer) IndexPublications(
	ctx context.Context,
	data chan crossrefindexer.SimplifiedPublication,
) error {
	// TODO: Add flag to delete existing index if wanted
	countSuccessful := &atomic.Uint64{}
	start := time.Now()
	bulkIndexer, err := createBulkIndexer(i.config, i.client)
	if err != nil {
		return err
	}

	for {
		// we receive a publication
		pub, stillOpen := <-data
		if !stillOpen {
			// If the channel is closed  - We "commit" the publications already in the slice before returning
			err := bulkIndexer.Close(ctx)
			i.logStats(bulkIndexer.Stats(), start)
			return errors.Wrap(err, "Closing of bulkindexer failed")
		}

		jsonData, err := json.Marshal(pub)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("Cannot encode publication %s", pub.DOI))
		}

		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		err = bulkIndexer.Add(
			ctx,
			i.bulkIndexerItem(bulkIndexer, pub.DOI, jsonData, countSuccessful, start),
		)
		if err != nil {
			return errors.Wrap(err, "Adding of indexing item failed")
		}
	}
}

// bulkIndexerItem builds the object to be passed for indexing
func (i *Indexer) bulkIndexerItem(
	bulkIndexer esutil.BulkIndexer,
	documentId string,
	data []byte,
	countSuccessful *atomic.Uint64,
	startTime time.Time,
) esutil.BulkIndexerItem {
	return esutil.BulkIndexerItem{
		Action:     "index",
		DocumentID: documentId,
		Body:       bytes.NewReader(data),

		// OnSuccess is called for each successful operation
		OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
			count := countSuccessful.Add(1)
			if count%10000 == 0 {
				i.logStats(bulkIndexer.Stats(), startTime)
			}
		},
		OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
			if err != nil {
				i.log.Errorw("Indexing failed", "err", err)
			} else {
				i.log.Errorw("Indexing failed", "type", res.Error.Type, "reason", res.Error.Reason)
			}
		},
	}
}

func (i *Indexer) logStats(biStats esutil.BulkIndexerStats, start time.Time) {
	dur := time.Since(start)

	if biStats.NumFailed > 0 {
		i.log.Errorf(
			"Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			humanize.Comma(int64(biStats.NumFailed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	} else {
		i.log.Infof(
			"Sucessfuly indexed [%s] documents in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	}
}

func createElasticClient(
	cfg Config,
	retryBackoff *backoff.ExponentialBackOff,
) (*elasticsearch.Client, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		RetryOnStatus:       []int{502, 503, 504, 429},
		Password:            cfg.Password,
		Username:            cfg.Username,
		Addresses:           cfg.Addresses,
		CACert:              cfg.CACert,
		DisableRetry:        cfg.DisableRetry,
		CompressRequestBody: cfg.CompressRequestBody,
		MaxRetries:          cfg.MaxRetries,
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}
			return retryBackoff.NextBackOff()
		},
	})

	return es, errors.Wrap(err, "failed to init elasticsearch client")
}

func createBulkIndexer(cfg Config, es *elasticsearch.Client) (esutil.BulkIndexer, error) {
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         cfg.IndexName,                                  // The default index name
		Client:        es,                                             // The Elasticsearch client
		NumWorkers:    cfg.NumWorkers,                                 // The number of worker goroutines
		FlushBytes:    cfg.FlushBytes,                                 // The flush threshold in bytes
		FlushInterval: time.Duration(cfg.FlushInterval) * time.Second, // The periodic flush
	})

	return bi, errors.Wrap(err, "could not create bulk indexer")
}
