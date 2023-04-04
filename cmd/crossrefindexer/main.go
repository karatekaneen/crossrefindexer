package main

import (
	"flag"
	"github.com/karatekaneen/crossrefindexer"
	"log"
	"os"
)

func main() {
	log.Println("hello")
	log.Println(os.Getwd())

	var dataPath string
	flag.StringVar(
		&dataPath,
		"path",
		os.Getenv("POOP"),
		"Path to the crossref data, can be both directory or a single file.",
	)
	flag.Parse()
	crossrefindexer.Load(dataPath)
}
