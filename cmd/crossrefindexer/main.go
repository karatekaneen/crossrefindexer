package main

import (
	"flag"
	"github.com/karatekaneen/crossrefindexer"
	"io"
	"log"
	"os"
	"path/filepath"
)

func ReadFile(path string) {
	file, err := os.Open(path)
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

	if dataPath == "" {
		log.Fatalf(dataPath)
	}

	err := filepath.Walk(dataPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}
		if info.IsDir() {
			return nil
		} else {
			ReadFile(path)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
}
