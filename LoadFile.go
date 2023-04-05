package crossrefindexer

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func readFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("file %q: %w", path, err)
	}
	defer file.Close()

	gzr, err := GzipReader(file)
	if err != nil {
		return fmt.Errorf("GzipReader: %w", err)
	}

	fileType, err := ClassifyDataFormat(gzr)
	if err != nil {
		return fmt.Errorf("Err with identifying json format: %w", err)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("Err with reseting to beiging of the buffer: %w", err)
	}

	gzr, err = GzipReader(file)
	if err != nil {
		return fmt.Errorf("GzipReader: %w", err)
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
		return fmt.Errorf("Err with parsing json file of type %q: %w", fileType, err)
	}

	log.Println("Found", len(pubs))
	return nil

}

func Load(path string) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty: %q", path)
	}

	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf(": %w", err)
		}
		if !info.IsDir() && filepath.Ext(path) == ".gz" {
			return readFile(path)
		}
		// if !info.IsDir() && filepath.Ext(path) == ".gz" {
		// 	publications := make(chan CrossRef)
		// 	go func() {
		// 		defer close(publications)
		// 		if err := readFile(path, publications); err != nil {
		// 			log.Printf("Err with reading file %q: %v", path, err)
		// 		}
		// 	}()
		// }
		return nil
	})
}
