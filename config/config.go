package config

import (
	"github.com/alecthomas/kong"
	"github.com/karatekaneen/crossrefindexer/elastic"
)

const description = `Small CLI application to uncompress and index Crossref metadata. 
It can read from file, directories and stdin.
It supports both compressed (xz only at the time of writing) and raw JSON/NDJSON.`

type Config struct {
	File     string         `help:"Absolute or relative path to a single file to index. If you set to '-' it will read from stdin"           short:"f" optional:"" type:"existingfile"`
	Dir      string         `help:"Absolute or relative path to a directory containing files to index"                                                 optional:"" type:"existingdir"`
	Elastic  elastic.Config `help:"Configuration for elasticsearch connection and indexing"                                                            optional:""                     embed:"" prefix:"es."`
	Format   string         `help:"The format of the uncompressed files. Will try to detect if not provided. Can be json, ndjson or unknown"           optional:""                                           default:"unknown" enum:"unknown,json,ndjson"`
	LogLevel string         `help:"Log verbosity. Can be debug, info, warn, error"                                                                                                                           default:"info"                               name:"loglevel"`
}

func Load() *Config {
	c := Config{}

	// TODO: Add so that config can be loaded from file
	ctx := kong.Parse(
		&c,
		kong.UsageOnError(),
		kong.Description(description),
	)

	if err := ctx.Validate(); err != nil {
		//nolint:errcheck
		ctx.PrintUsage(false)
		ctx.Fatalf("config validation failed: %v", err)
	} else if c.Dir == "" && c.File == "" {
		//nolint:errcheck
		ctx.PrintUsage(false)
		ctx.Fatalf("Either dir or file must be provided")
	}

	return &c
}
