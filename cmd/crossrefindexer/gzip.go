package main

import (
	"compress/gzip"
	_ "encoding/json"
	_ "log"
	"os"
)

// GzipReader creates a new Reader of a gziped file.
//
// It is the caller's responsibility to call Close on the Reader when done.
func GzipReader(source string) (*gzip.Reader, error) {
	file, err := os.Open(source)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}

	return gzipReader, nil
}

// func main() {
// 	r, _ := GzipReader("D1000000.json.gz")
// 	j := json.NewDecoder(r)
// 	var elm map[string]any
// 	err := j.Decode(&elm)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	log.Println(elm)

// }
