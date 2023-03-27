package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

// ClassifyDataFormat ...
func ClassifyDataFormat(r io.ReadSeeker) string {
	d := json.NewDecoder(r)
	token, _ := d.Token()
	token, _ = d.Token()
	defer r.Seek(0, 0)
	if token == "items" {
		println(token, d.InputOffset(), "crossref new")
		return "json"
	}
	println(token, d.InputOffset(), "crossref old")
	return "jsonl"
}

func main() {
	file, err := os.Open("testdata/2022/0.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ClassifyDataFormat(file)

	d := json.NewDecoder(file)
	token, _ := d.Token()
	token, _ = d.Token()
	println(token, d.InputOffset())
}
