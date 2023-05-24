package config

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/karatekaneen/crossrefindexer"
	"github.com/karatekaneen/crossrefindexer/elastic"
)

const description = `Small CLI application to uncompress and index Crossref metadata. 
It can read from file, directories and stdin.
It supports both compressed (gzip only at the time of writing) and raw JSON/NDJSON.`

type Config struct {
	RemoveIndex bool                   `help:"Remove existing index before starting. WARNING - you will not get any confirmation prompt"                                                            default:"false"`
	File        string                 `help:"Absolute or relative path to a single file to index. If you set to '-' it will read from stdin"                                                                         short:"f" optional:"" type:"existingfile"`
	Dir         string                 `help:"Absolute or relative path to a directory containing files to index"                                                                                                               optional:"" type:"existingdir"`
	Elastic     elastic.Config         `help:"Configuration for elasticsearch connection and indexing"                                                                                                                          optional:""                     embed:"" prefix:"es."`
	Format      crossrefindexer.Format `help:"The format of the uncompressed files. Will try to detect if not provided but is required if using stdin. Can be json, ndjson or unknown"              default:"unknown"           optional:""                                           enum:"unknown,json,ndjson"`
	Compression string                 `help:"How the data file is compressed. For files it will use the file extension if not provided. For dirs it will be ignored. Can be unknown, none or gzip" default:"unknown" short:"c"                                                       enum:"unknown,none,gzip"`
	LogLevel    string                 `help:"Log verbosity. Can be debug, info, warn, error"                                                                                                       default:"info"                                                                                               name:"loglevel"`
}

type configValidator func(Config) error

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
	}

	for _, validator := range []configValidator{hasPath, hasFormat, hasCompression} {
		if err := validator(c); err != nil {
			//nolint:errcheck
			ctx.PrintUsage(false)
			ctx.Fatalf("config validation failed: %v", err)
		}
	}

	return &c
}

func hasPath(c Config) error {
	if c.Dir == "" && c.File == "" {
		return fmt.Errorf("Either dir or file must be provided")
	}
	return nil
}

func hasFormat(c Config) error {
	if c.Format == "unknown" && c.File == "-" {
		return fmt.Errorf("Format must be specified when reading from stdin")
	}
	return nil
}

func hasCompression(c Config) error {
	if c.Compression == "unknown" && c.File == "-" {
		return fmt.Errorf("Compression must be specified when reading from stdin")
	}
	return nil
}
