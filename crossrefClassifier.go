package main

import (
	"encoding/json"
	"log"
	"os"
)

func CrossrefVerClassifier(d *json.Decoder) int {
	token, _ := d.Token()
	token, _ = d.Token()
	if token == "items" {
		println(token, d.InputOffset(), "crossref new")
		return 1
	}
	println(token, d.InputOffset(), "crossref old")
	return 2
}

func main() {

	file, err := os.Open("testdata/gap/D1000000.json")
	// file, err := os.Open("testdata/2022/0.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	d := json.NewDecoder(file)

	CrossrefVerClassifier(d)

}
