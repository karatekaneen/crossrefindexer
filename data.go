package crossrefindexer

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Format string

const (
	FormatUnknown Format = "unknown"
	FormatJSON    Format = "json"
	FormatNDJSON  Format = "ndjson"
)

type DataContainer struct {
	Data        io.Reader // The data to index if passed by stdin or similar
	Path        string    // Path to the file to read
	Format      Format    // Format of the file, either "json" or "ndjson"
	Compression string    // The kind of compression. Currently only supports "none" or "gzip"
}

func (d *DataContainer) Valid() error {
	if d.Data == nil || d.Format == "unknown" || d.Compression == "unknown" {
		return fmt.Errorf("DataContainer invalid: %+v", d)
	}
	return nil
}

// readJsonData consumes the reader of uncompressed data.
// Supports both regular json and newline delimited json (ndjson)
func readJsonData(r io.Reader, ch chan Crossref, format Format) error {
	d := json.NewDecoder(r)

	// The json format is quite nested so we need to skip
	// three levels "deep" to reach the data that we want
	if format == "json" {
		//nolint:errcheck
		d.Token()
		//nolint:errcheck
		d.Token()
		//nolint:errcheck
		d.Token()
	}

	elementIndex := 0

	for d.More() {
		var publication Crossref

		err := d.Decode(&publication)
		if err == io.EOF {
			break
		} else if err != nil {
			return errors.Wrapf(err, "failed on parsing element %d", elementIndex)
		}

		ch <- publication
		elementIndex++
	}
	return nil
}

// ParseData reads the data described in the container and passes it via the out channel
func ParseData(container DataContainer, out chan Crossref) error {
	// Declare the variables so that we don't shadow them
	var (
		rawData, data io.ReadCloser
		err           error
	)

	if container.Data != nil {
		rawData = io.NopCloser(container.Data)
	} else {
		rawData, err = os.Open(container.Path)
		if err != nil {
			return fmt.Errorf("could not open file: %w", err)
		}
	}
	defer rawData.Close() // Make sure we close before we return

	data = rawData
	if container.Compression == "gzip" {
		data, err = gzip.NewReader(data)
		if err != nil {
			return fmt.Errorf("create gzip reader: %w", err)
		}
	}
	defer data.Close() // Close the gzipped data as well.

	if err := readJsonData(data, out, container.Format); err != nil {
		return fmt.Errorf(
			"err with parsing data of type %q and format %q: %w",
			container.Format,
			container.Format,
			err,
		)
	}

	return nil
}

// Load structures the data that should be indexed.
// It returns a slice of items to be processed.
func Load(
	logger *zap.SugaredLogger,
	path, dir string, format Format, compression string,
	data io.Reader,
) ([]DataContainer, error) {
	switch {
	case path == "-":
		logger.Debugw("Reading from stdin")
		return []DataContainer{
			{
				Data:        data,
				Format:      format,
				Compression: compression,
			},
		}, nil
	case path != "":
		logger.Debug("Reading from file")
		d, err := dataContainerFromPath(path, format, compression)
		return []DataContainer{d}, err
	case dir != "":
		filesInDir, err := listFiles(dir)
		if err != nil {
			return nil, fmt.Errorf("File listing failed: %w", err)
		}

		output := make([]DataContainer, 0, len(filesInDir))

		for _, f := range filesInDir {
			d, err := dataContainerFromPath(f, format, compression)
			if err != nil {
				return nil, fmt.Errorf("could not create datacontainer from path %s, : %w", f, err)
			}

			output = append(output, d)
		}

		return output, nil
	default:
		return nil, fmt.Errorf("Unexpected path")
	}
}

func listFiles(root string) ([]string, error) {
	files := []string{}

	acceptedExtensions := map[string]struct{}{
		".gzip":   {},
		".ndjson": {},
		".json":   {},
		".gz":     {},
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if _, accepted := acceptedExtensions[filepath.Ext(path)]; !info.IsDir() && accepted {
			files = append(files, path)
		}

		return nil
	})
	return files, err
}

func dataContainerFromPath(path string, format Format, compression string) (DataContainer, error) {
	d := DataContainer{
		Format:      format,
		Compression: compression,
		Path:        path,
	}

	// Don't override if compression has been set explicitly
	if d.Compression == "unknown" || d.Compression == "" {
		ext := filepath.Ext(path)
		switch ext {
		case ".gzip":
			fallthrough
		case ".gz":
			d.Compression = "gzip"
		case ".ndjson":
			fallthrough
		case ".json":
			d.Compression = "none"
		default:
			d.Compression = "none"
		}
	}

	if d.Format == "unknown" {
		// ClassifyFormat()
		// TODO: Add classification here
		log.Fatalln("Format detection not implemented")
	}

	return d, nil
}

// ClassifyDataFormat Tries to figure out if the format is JSON or JSONL/NDJson (Newline Delimited JSON)
// func ClassifyDataFormat(r io.Reader) (string, error) {
// 	d := json.NewDecoder(r)
//
// 	_, err := d.Token()
// 	if err != nil {
// 		return "", err
// 	}
//
// 	token, err := d.Token()
// 	if err != nil {
// 		return "", err
// 	}
//
// 	dataFormat := "jsonl"
// 	if token == "items" {
// 		dataFormat = "json"
// 	}
//
// 	return dataFormat, nil
// }
