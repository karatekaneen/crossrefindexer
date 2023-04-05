package main

import (
	"flag"
	"log"
	"os"

	"github.com/karatekaneen/crossrefindexer"
)

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
	// TODO: Add indexing around here

	go func() {
		err := crossrefindexer.Load("testdata/2021/", publications)
		if err != nil {
			log.Fatal(err)
		}
	}()

	for {
		pub, open := <-publications
		if !open {
			break
		}
		log.Println(pub.Doi)
	}
}
