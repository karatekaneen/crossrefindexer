package crossrefindexer

import (
	"encoding/json"
	"io"
)

// ClassifyDataFormat ...
func ClassifyDataFormat(r io.Reader) (string, error) {
	d := json.NewDecoder(r)

	_, err := d.Token()
	if err != nil {
		return "", err
	}

	token, err := d.Token()
	if err != nil {
		return "", err
	}

	dataFormat := "jsonl"
	if token == "items" {
		dataFormat = "json"
	}

	return dataFormat, nil
}
