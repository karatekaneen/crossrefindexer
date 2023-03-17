# crossrefindexer

Indexes metadata from Crossref into Elasticsearch. Primarily to be used with [Biblio-Glutton](https://github.com/kermitt2/biblio-glutton)

## TODO

### CLI

- `crossrefindexer -path /path/to/directory/`
- `crossrefindexer -path /path/to/directory/ -format jsonl`
- `crossrefindexer -path /path/to/data/0.json.gz -url http://localhost:9200`

#### Flags

- reset: should delete the index before running

### GZIP (step #1)

- Unzip expose as a `io.Reader`

### JSON parsing (Step #2)

- Be able to tell what format we get as input - (`items` array OR jsonl (Step #3) format from the gap folder)
- Convert into a specified struct `CrossrefPublication`: https://mholt.github.io/json-to-go/
- Takes the io.Reader as input
- Parses each single item
- Posts each entry on a channel (Step #4)

### Convert to indexed format

- Takes a single `CrossrefPublication` returns `SimplifiedPublication`. See [these functions](https://github.com/kermitt2/biblio-glutton/blob/master/indexing/main.js#L217-L366) for expected format
- Sends the converted entry to another channel

### Bulk batching into elasticsearch

- Receives entries (`SimplifiedPublication`) on channel. Indexes when it has received 1000 OR that the channel is closed
