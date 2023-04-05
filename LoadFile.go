package crossrefindexer

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func readFile(path string, publications chan CrossRef) error {
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
		return fmt.Errorf("err with identifying json format: %w", err)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("err with reseting to beiging of the buffer: %w", err)
	}

	gzr, err = GzipReader(file)
	if err != nil {
		return fmt.Errorf("GzipReader: %w", err)
	}

	if err := JsonParser(gzr, publications, fileType); err != nil {
		return fmt.Errorf("err with parsing json file of type %q: %w", fileType, err)
	}

	return nil

}

func Load(path string, publications chan CrossRef) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty: %q", path)
	}

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf(": %w", err)
		}

		if !info.IsDir() && filepath.Ext(path) == ".gz" {
			if err := readFile(path, publications); err != nil {
				log.Printf("Err with reading file %q: %v", path, err)
			}
		}
		return nil
	})

	close(publications)
	return err
}
