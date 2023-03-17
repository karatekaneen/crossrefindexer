package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

func main() {
	log.Println("hello")

	file, err := os.Open("testdata/2022/0.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	r := bufio.NewReader(file)
	d := json.NewDecoder(r)

	i := 0

	// t, err := d.Token()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// TODO: Figure out how this works:
	d.Token()
	d.Token()
	d.Token()
	// log.Println("token", t)
	for d.More() {
		var elm map[string]any
		err := d.Decode(&elm)
		if err != nil {
			log.Println(err)
		}

		// log.Println("*****", i, "******")
		// for key, value := range elm {
		// 	log.Println(key, ":", value)
		// }
		i++

		// if i > 5 {
		// 	break
		// }
	}
	d.Token()
	log.Println(i)
}
