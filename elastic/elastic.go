package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

type Config struct {
	IndexName           string        `help:"The index to write to"                                    default:"crossref"              name:"index"         env:"ES_INDEX"`
	FlushBytes          int           `help:"How many bytes to buffer before flushing. Defaults to 5M" default:"5000000"               name:"flushbytes"    env:"ES_FLUSH_BYTES"`
	FlushInterval       time.Duration `help:"How many seconds to wait before flushing"                 default:"10s"                   name:"flushinterval" env:"ES_FLUSH_INTERVAL"`
	NumWorkers          int           `help:"Number of goroutines to run"                              default:"2"                     name:"workers"       env:"ES_WORKERS"`
	Password            string        `help:"Password to elasticsearch"                                                                                     env:"ES_PASSWORD"       short:"p" optional:""`
	Username            string        `help:"Username to elasticsearch"                                                                                     env:"ES_USER"           short:"u" optional:""`
	Addresses           []string      `help:"Elasticsearch hosts"                                      default:"http://127.0.0.1:9200" name:"hosts"         env:"ES_HOSTS"`
	CACert              []byte        `help:"CA cert to trust"                                                                         name:"ca"            env:"ES_CA_CERT"                  optional:""`
	DisableRetry        bool          `help:"Fail on first failure"                                    default:"false"                 name:"noretry"       env:"ES_NO_RETRY"`
	MaxRetries          int           `help:"Max number of retries after failure"                      default:"5"                                          env:"ES_MAX_RETRIES"`
	CompressRequestBody bool          `help:"If the request body should be compressed"                 default:"false"                 name:"compress"      env:"ES_COMPRESS"`
}

type Indexer struct {
	config    Config
	client    *elasticsearch.Client
	log       *zap.SugaredLogger
	transport http.RoundTripper
}

type Option func(*Indexer)

func New(config Config, log *zap.SugaredLogger, options ...Option) (*Indexer, error) {
	idx := &Indexer{
		config: config,
		log:    log,
	}
	// Instantiate the exponential backoff thingy
	retryBackoff := backoff.NewExponentialBackOff()

	for _, option := range options {
		option(idx)
	}

	esClient, err := createElasticClient(config, retryBackoff, idx.transport, log)
	if err != nil {
		return nil, err
	}

	idx.client = esClient

	return idx, nil
}

// Add custom transport.
// TODO: Maybe move this to config?
func WithTransport(rt http.RoundTripper) Option { return func(i *Indexer) { i.transport = rt } }

func (i *Indexer) DeleteIndex(ctx context.Context, indexName string) error {
	// The API is kinda fubar so lets just assign it to a variable for ease of use
	deleteApi := i.client.API.Indices.Delete

	resp, err := deleteApi([]string{indexName})
	if err != nil {
		return fmt.Errorf("Delete request failed: %w", err)
	}
	defer resp.Body.Close()

	// Let's ignore 404s
	if resp.IsError() && resp.StatusCode != http.StatusNotFound {
		return newElasticError(resp)
	}

	return nil
}

func (i *Indexer) CreateIndex(ctx context.Context, indexName string, settings IndexSettings) error {
	createApi := i.client.API.Indices.Create

	data, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("could not marshal settings to json: %w", err)
	}
	resp, err := createApi(
		indexName,
		createApi.WithContext(ctx),
		createApi.WithBody(bytes.NewReader(data)),
	)
	if err != nil {
		return fmt.Errorf("Delete request failed: %w", err)
	}
	defer resp.Body.Close()

	// Let's ignore conflicts where the index already exists
	if resp.IsError() {
		err := newElasticError(resp)
		if !strings.Contains(err.Error(), "already exists") {
			return err
		}
	}
	return nil
}

// IndexPublications is responsible for consuming all the publications sent on `data` and then close
// the indexer when the channel is closed
func (i *Indexer) IndexPublications(
	ctx context.Context,
	data chan crossrefindexer.SimplifiedPublication,
) error {
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
			i.log.Debug("Starting to close indexer")
			err := bulkIndexer.Close(ctx)
			i.log.Debugf("Closed indexer with err: %q", err)

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

			i.log.Debugln(count)
			// Log more often in the beginning to get quick feedback
			highFreq := count < 1_000_000 && count%100_000 == 0   // Log every 100k in the beginning
			lowFreq := count >= 1_000_000 && count%1_000_000 == 0 // Log every 1m afterwards
			if highFreq || lowFreq {
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
	transport http.RoundTripper,
	logger *zap.SugaredLogger,
) (*elasticsearch.Client, error) {
	elasticConfig := elasticsearch.Config{
		RetryOnStatus:       []int{502, 503, 504, 429},
		Password:            cfg.Password,
		Username:            cfg.Username,
		Addresses:           cfg.Addresses,
		CACert:              cfg.CACert,
		DisableRetry:        cfg.DisableRetry,
		CompressRequestBody: cfg.CompressRequestBody,
		MaxRetries:          cfg.MaxRetries,
		Transport:           transport,
		RetryBackoff: func(i int) time.Duration {
			if i == 1 {
				retryBackoff.Reset()
			}

			logger.Debugf("Retry for the %d time", i)
			return retryBackoff.NextBackOff()
		},
	}

	es, err := elasticsearch.NewClient(elasticConfig)

	return es, errors.Wrap(err, "failed to init elasticsearch client")
}

func createBulkIndexer(cfg Config, es *elasticsearch.Client) (esutil.BulkIndexer, error) {
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         cfg.IndexName,     // The default index name
		Client:        es,                // The Elasticsearch client
		NumWorkers:    cfg.NumWorkers,    // The number of worker goroutines
		FlushBytes:    cfg.FlushBytes,    // The flush threshold in bytes
		FlushInterval: cfg.FlushInterval, // The periodic flush
		Refresh:       "-1",              // Disable refresh
	})

	return bi, errors.Wrap(err, "could not create bulk indexer")
}
