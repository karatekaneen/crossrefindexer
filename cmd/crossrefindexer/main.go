package main

import (
	"context"
	"log"
	"sync"

	"github.com/karatekaneen/crossrefindexer"
	"github.com/karatekaneen/crossrefindexer/config"
	"github.com/karatekaneen/crossrefindexer/elastic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// type indexer interface {
// 	Index(ctx context.Context, data chan crossrefindexer.SimplifiedPublication) error
// }

func createLogger(level string) (*zap.Logger, error) {
	var l zapcore.Level

	switch level {
	case "debug":
		l = zap.DebugLevel
	case "info":
		l = zap.InfoLevel
	case "warn":
		l = zap.WarnLevel
	case "error":
		l = zap.ErrorLevel
	default:
		l = zap.InfoLevel
	}

	loggerSettings := zap.NewDevelopmentConfig()
	loggerSettings.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	loggerSettings.Level = zap.NewAtomicLevelAt(l)

	return loggerSettings.Build()
}

func main() {
	cfg := config.Load()
	// Init logger
	l, err := createLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal(err)
	}
	logger := l.Sugar().With(
		"file", cfg.File,
		"dir", cfg.Dir,
		"format", cfg.Format,
		"compression", cfg.Compression,
	)
	logger.Debugln("Config loaded successfully")

	publications := make(chan crossrefindexer.Crossref)
	dataToIndex := make(chan crossrefindexer.SimplifiedPublication)

	// LoadData. Can be file (json/gzip), dir or stdin
	// If file: get format & compression then read data
	// If dir: walk files, extract format, infer compression and then read as file

	// go func() {
	// 	err := crossrefindexer.Load("testdata/2021/0.json.gz", publications)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()

	// cfg := elastic.Config{
	// 	IndexName:           "wtf",
	// 	Addresses:           []string{"http://localhost:9200"},
	// 	Username:            "elastic",
	// 	Password:            "123change...",
	// 	CompressRequestBody: true,
	// 	MaxRetries:          5,
	// 	NumWorkers:          4,
	// }
	//
	es, err := elastic.New(cfg.Elastic, logger)
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
