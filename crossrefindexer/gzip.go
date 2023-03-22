package main

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"log"
	"os"
)

// GzipReader creates a new Reader of a gziped file.
//
// It is the caller's responsibility to call Close on the Reader when done.
func GzipReader(r io.Reader) (*gzip.Reader, error) {
	gzipReader, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return gzipReader, nil
}

func main() {
	file, err := os.Open("testdata/gap/D1000000.json.gz")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	r, err := GzipReader(file)
	if err != nil {
		log.Fatal(err)
	}

	j := json.NewDecoder(r)
	var elm map[string]any
	err = j.Decode(&elm)
	if err != nil {
		log.Println(err)
	}
	log.Println(elm)

}
