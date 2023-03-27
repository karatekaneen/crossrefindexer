package main

import (
	"encoding/json"
	"log"
	"os"
)

// ClassifyDataFormat ...
func ClassifyDataFormat(f string) string {
	file, err := os.Open(f)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	d := json.NewDecoder(file)
	token, _ := d.Token()
	token, _ = d.Token()
	if token == "items" {
		return "json"
	}
	return "jsonl"
}

func main() {
	jsonType := ClassifyDataFormat("testdata/gap/D1000000.json")
	if jsonType == "json" {
		println("crossref new")
	} else {
		println("crossref old")
	}
	// 	file, err := os.Open("testdata/gap/D1000000.json")
	// 	// file, err := os.Open("testdata/2022/0.json")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	defer file.Close()

	// 	d := json.NewDecoder(file)
	// 	ClassifyDataFormat(d)
	// 	file.Seek(0, 0)
	// 	ClassifyDataFormat(d)

	// 	c := json.NewDecoder(file)
	// 	ClassifyDataFormat(c)

}
