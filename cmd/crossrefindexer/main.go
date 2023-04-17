package main

import (
	"context"
	"flag"
	_ "fmt"
	"github.com/karatekaneen/crossrefindexer"
	"log"
	"os"
)

type indexer interface {
	Index(ctx context.Context, data chan crossrefindexer.SimplifiedPublication) error
}

func main() {
	log.Println("hello")
	log.Println(os.Getwd())

	var dataPath string
	flag.StringVar(
		&dataPath,
		"path",
		os.Getenv("POOP"),
		"Path to the crossref data, can be both directory or a single file.",
	)
	flag.Parse()

	publications := make(chan crossrefindexer.CrossRef)

	// TODO: Add conversion
	// TODO: Add indexing around here `indexer.Index(ctx, convertedPublication)`

	go func() {
		err := crossrefindexer.Load("testdata/2021/0.json.gz", publications)
		if err != nil {
			log.Fatal(err)
		}
	}()

	for {
		pub, open := <-publications
		if !open {
			break
		}
		test := crossrefindexer.ToSimplifiedPublication(&pub)
		log.Printf("%v", test.Bibliographic)
	}
}
