package main

import (
	"context"
	"flag"
	"log"
	"os"
	"sync"

	"github.com/karatekaneen/crossrefindexer"
	"github.com/karatekaneen/crossrefindexer/elastic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// type indexer interface {
// 	Index(ctx context.Context, data chan crossrefindexer.SimplifiedPublication) error
// }

func createLogger() (*zap.Logger, error) {
	// var loggerSettings zap.Config

	// if cfg.Env == "production" {
	// 	loggerSettings = zap.NewProductionConfig()
	// 	loggerSettings.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	// } else {
	loggerSettings := zap.NewDevelopmentConfig()
	loggerSettings.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	// }

	return loggerSettings.Build()
}

func main() {
	l, err := createLogger()
	if err != nil {
		log.Fatal(err)
	}

	logger := l.Sugar()
	var dataPath string
	flag.StringVar(
		&dataPath,
		"path",
		os.Getenv("POOP"),
		"Path to the crossref data, can be both directory or a single file.",
	)
	flag.Parse()

	publications := make(chan crossrefindexer.Crossref)
	dataToIndex := make(chan crossrefindexer.SimplifiedPublication)

	// TODO: Add conversion
	// TODO: Add indexing around here `indexer.Index(ctx, convertedPublication)`

	go func() {
		err := crossrefindexer.Load("testdata/2021/0.json.gz", publications)
		if err != nil {
			log.Fatal(err)
		}
	}()

	cfg := elastic.Config{
		IndexName:           "wtf",
		Addresses:           []string{"http://localhost:9200"},
		Username:            "elastic",
		Password:            "123change...",
		CompressRequestBody: true,
		MaxRetries:          5,
		NumWorkers:          4,
	}

	es, err := elastic.New(cfg, logger)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		if err := es.IndexPublications(context.Background(), dataToIndex); err != nil {
			log.Fatal(err)
		}
		log.Println(err)
		wg.Done()
	}()

	count := 0
	for {
		pub, open := <-publications
		if !open {
			close(dataToIndex)
			break
		}
		count++
		dataToIndex <- crossrefindexer.ToSimplifiedPublication(&pub)
	}
	log.Println("count", count)
	wg.Wait()
}
