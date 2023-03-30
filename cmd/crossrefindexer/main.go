package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/karatekaneen/crossrefindexer"
)

func main() {
	log.Println("hello")
	var dataPath string
	flag.StringVar(
		&dataPath,
		"path",
		os.Getenv("POOP"),
		"Path to the crossref data, can be both directory or a single file.",
	)
	flag.Parse()

	if dataPath != "" || dataPath == "" {
		log.Fatalf(dataPath)
	}

	// TODO: If path is directory: Read all files in directory
	// TODO: Else just read the single file

	file, err := os.Open("testdata/2022/0.json.gz")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	gzr, err := crossrefindexer.GzipReader(file)
	if err != nil {
		log.Fatal(err)
	}

	fileType, err := crossrefindexer.ClassifyDataFormat(gzr)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		log.Fatal(err)
	}

	gzr, err = crossrefindexer.GzipReader(file)
	if err != nil {
		log.Fatal(err)
	}

	publications := make(chan crossrefindexer.CrossRef)

	pubs := []crossrefindexer.CrossRef{}

	go func() {
		for {
			pub, open := <-publications
			if !open {
				break
			}
			// TODO: Convert to simplified format
			// TODO: Send on another channel
			pubs = append(pubs, pub)
		}
	}()

	if err := crossrefindexer.JsonParser(gzr, publications, fileType); err != nil {
		log.Fatal(err)
	}

	log.Println("Found", len(pubs))

}
