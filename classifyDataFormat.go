package crossrefindexer

import (
	"encoding/json"
	"io"
)

// ClassifyDataFormat ...
func ClassifyDataFormat(r io.ReadSeeker) (string, error) {
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
	_, err = r.Seek(0, 0)
	if err != nil {
		return "", err
	}

	return dataFormat, nil
}
