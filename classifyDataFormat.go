package crossrefindexer

import (
	"encoding/json"
	"io"
	"log"
)

// ClassifyDataFormat ...
func ClassifyDataFormat(r io.ReadSeeker) string {
	d := json.NewDecoder(r)

	token, err := d.Token()
	if err != nil {
		log.Fatal("error geting the first token\n", err)
	}

	token, err = d.Token()
	if err != nil {
		log.Fatal("error geting the second token\n", err)
	}

	dataFormat := ""
	if token == "items" {
		dataFormat = "json"
	} else {
		dataFormat = "jsonl"
	}

	loc, err := r.Seek(0, 0)
	if err != nil {
		log.Fatal(loc, err)
	}

	return dataFormat
}
