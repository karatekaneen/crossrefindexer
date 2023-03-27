package main

import (
	"encoding/json"
	"fmt"
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
		return "json"
	}
	return "jsonl"
}

func main() {
	file, err := os.Open("testdata/gap/D1000000.json")
	// file, err := os.Open("testdata/2022/0.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	jsonType := ClassifyDataFormat(file)
	fmt.Printf("Json Format Type:%v\n", jsonType)
	d := json.NewDecoder(file)
	token, _ := d.Token()
	token, _ = d.Token()
	fmt.Printf("Token:%v\nDecoder buf offset:%v\n", token, d.InputOffset())
}
