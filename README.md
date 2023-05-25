# crossrefindexer

Indexes metadata from Crossref into Elasticsearch. Primarily to be used with [Biblio-Glutton](https://github.com/kermitt2/biblio-glutton).
It is currently a single-purpose application that has the format for Glutton hardcoded.
If you want to modify what data is being indexed you need to modify
the `ToSimplifiedPublication` function in the root package.

This application can read both regular JSON as well as newline-delimited JSON (NDJSON).
It supports GZIP and uncompressed data (TAR will be added).
You can read from single files, directories or stdin.
Configuration can be done via commandline flags or env variables.

## Installation

Make sure you have Go installed and made sure that the correct folders are added to $PATH.
Then run:

```sh
go install github.com/karatekaneen/crossrefindexer/cmd/crossrefindexer`
```

## Usage

### Configuration

To see full configuration options run `crossrefindexer --help`.
The output below is how it looks at the time of writing.

```sh
Usage: crossrefindexer

Small CLI application to uncompress and index Crossref metadata. It can read
from file, directories and stdin. It supports both compressed (gzip only at the
time of writing) and raw JSON/NDJSON.

Flags:
  -h, --help                     Show context-sensitive help.
      --remove-index             Remove existing index before starting. WARNING
                                 - you will not get any confirmation prompt
  -f, --file=STRING              Absolute or relative path to a single file
                                 to index. If you set to '-' it will read from
                                 stdin
      --dir=STRING               Absolute or relative path to a directory
                                 containing files to index
      --es.index="crossref"      The index to write to ($ES_INDEX)
      --es.flushbytes=5000000    How many bytes to buffer before flushing.
                                 Defaults to 5M ($ES_FLUSH_BYTES)
      --es.flushinterval=10s     How many seconds to wait before flushing
                                 ($ES_FLUSH_INTERVAL)
      --es.workers=4             Number of goroutines to run ($ES_WORKERS)
  -p, --es.password=STRING       Password to elasticsearch ($ES_PASSWORD)
  -u, --es.username=STRING       Username to elasticsearch ($ES_USER)
      --es.hosts=http://127.0.0.1:9200,...
                                 Elasticsearch hosts ($ES_HOSTS)
      --es.ca=ES.CA,...          CA cert to trust ($ES_CA_CERT)
      --es.noretry               Fail on first failure ($ES_NO_RETRY)
      --es.max-retries=5         Max number of retries after failure
                                 ($ES_MAX_RETRIES)
      --es.compress              If the request body should be compressed
                                 ($ES_COMPRESS)
      --format="unknown"         The format of the uncompressed files. Will try
                                 to detect if not provided but is required if
                                 using stdin. Can be json, ndjson or unknown
  -c, --compression="unknown"    How the data file is compressed. For files it
                                 will use the file extension if not provided.
                                 For dirs it will be ignored. Can be unknown,
                                 none or gzip
      --loglevel="info"          Log verbosity. Can be debug, info, warn, error
```

### Read from stdin

When reading from stdin you must specify both format and compression.

```sh
# the part with "-f -" means that it is reading from stdin
cat testdata/2022/0.json.gz | crossrefindexer -f - --format json -c gzip
```

### Read from single file

```sh
# Compression is detected from the file extension
crossrefindexer -f testdata/2022/0.json.gz --format json
```

### Read from directory

```sh
# Compression is detected from the file extension.
# It supports multiple formats in the same directory.
crossrefindexer --dir testdata/2022 --format json
```

## TODO

- Support TAR files
