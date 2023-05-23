package main

import (
	"context"
	"log"
	"os"

	"github.com/karatekaneen/crossrefindexer"
	"github.com/karatekaneen/crossrefindexer/config"
	"github.com/karatekaneen/crossrefindexer/elastic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
)

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

	es, err := elastic.New(cfg.Elastic, logger)
	if err != nil {
		log.Fatal(err)
	}

	publications := make(chan crossrefindexer.Crossref)
	dataToIndex := make(chan crossrefindexer.SimplifiedPublication)

	// LoadData. Can be file (json/gzip), dir or stdin
	// If file: get format & compression then read data
	// If dir: walk files, extract format, infer compression and then read as file
	inputs, err := crossrefindexer.Load(
		logger,
		cfg.File,
		cfg.Dir,
		cfg.Format,
		cfg.Compression,
		os.Stdin,
	)
	if err != nil {
		logger.Fatalln(err)
	}

	group := new(errgroup.Group)               // Create an errgroup to manage goroutines
	indexGroup := new(errgroup.Group)          // A group to hold the indexing
	readGroup := new(errgroup.Group)           // A group to read the files
	readGroup.SetLimit(cfg.Elastic.NumWorkers) // Limit the number of files to read concurrently
	// TODO: Add flag to delete existing index if wanted

	// Initialize the indexing
	group.Go(func() error {
		indexGroup.Go(func() error {
			return es.IndexPublications(context.Background(), dataToIndex)
		})
		return indexGroup.Wait()
	})

	group.Go(func() error {
		defer close(publications)

		for index, container := range inputs {
			container := container // Because Go is wonky

			// Log progress
			logger.Debugw("Parsing file",
				"index", index,
				"numberOfItems", len(inputs),
				"path", container.Path,
			)

			// Process the file
			readGroup.Go(func() error { return crossrefindexer.ParseData(container, publications) })
		}

		return readGroup.Wait()
	})

	// Convert the data and pipe it to the indexing channel
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

	if err := group.Wait(); err != nil {
		logger.Fatalf("Something failed: %w", err)
	}

	logger.Infof("Indexed %d publications from %d files successfully", count, len(inputs))
}
