package crossrefindexer

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

func readFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := GzipReader(file)
	if err != nil {
		return err
	}

	fileType, err := ClassifyDataFormat(gzr)
	if err != nil {
		return err
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return err
	}

	gzr, err = GzipReader(file)
	if err != nil {
		return err
	}

	publications := make(chan CrossRef)

	pubs := []CrossRef{}

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

	if err := JsonParser(gzr, publications, fileType); err != nil {
		return err
	}

	log.Println("Found", len(pubs))
	return nil

}

func Load(path string) error {
	if path == "" {
		log.Fatalf(path)
	}

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println(err)
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".gz" {
			readFile(path)
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
