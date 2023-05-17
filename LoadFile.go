package crossrefindexer

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

type DataContainer struct {
	data        io.ReadCloser
	path        string
	format      string
	compression string
}

func (d *DataContainer) Valid() error {
	if d.data == nil || d.format == "unknown" || d.compression == "unknown" {
		return fmt.Errorf("DataContainer invalid: %+v", d)
	}
	return nil
}

func readFile(path string, publications chan Crossref) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("file %q: %w", path, err)
	}
	defer file.Close()

	gzr, err := GzipReader(file)
	if err != nil {
		return fmt.Errorf("GzipReader: %w", err)
	}

	fileType, err := ClassifyDataFormat(gzr)
	if err != nil {
		return fmt.Errorf("err with identifying json format: %w", err)
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("err with reseting to beiging of the buffer: %w", err)
	}

	gzr, err = GzipReader(file)
	if err != nil {
		return fmt.Errorf("GzipReader: %w", err)
	}

	if err := JsonParser(gzr, publications, fileType); err != nil {
		return fmt.Errorf("err with parsing json file of type %q: %w", fileType, err)
	}

	return nil
}

func Load(
	logger *zap.SugaredLogger,
	path, dir, format, compression string,
	data io.Reader,
) ([]DataContainer, error) {
	switch {
	case path == "-":
		logger.Debugw("Reading from stdin")
		return []DataContainer{
			{
				data:        io.NopCloser(data),
				format:      format,
				compression: compression,
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

func dataContainerFromPath(path, format, compression string) (DataContainer, error) {
	d := DataContainer{
		format:      format,
		compression: compression,
		path:        path,
	}

	// Don't override if compression has been set explicitly
	if d.compression == "unknown" || d.compression == "" {
		ext := filepath.Ext(path)
		switch ext {
		case ".gzip":
			fallthrough
		case ".gz":
			d.compression = "gzip"
		case ".ndjson":
			fallthrough
		case ".json":
			d.compression = "none"
		default:
			d.compression = "none"
		}
	}

	if d.format == "unknown" {
		// ClassifyFormat()
		// TODO: Add classification here
		log.Fatalln("Format detection not implemented")
	}

	return d, nil
}

// GzipReader creates a new Reader of a gziped file.
//
// It is the caller's responsibility to call Close on the Reader when done.
func GzipReader(r io.Reader) (*gzip.Reader, error) {
	gzipReader, err := gzip.NewReader(r)
	if err != nil {
		return gzipReader, fmt.Errorf("create gzip reader: %w", err)
	}
	return gzipReader, nil
}

// ClassifyDataFormat Tries to figure out if the format is JSON or JSONL/NDJson (Newline Delimited JSON)
func ClassifyDataFormat(r io.Reader) (string, error) {
	d := json.NewDecoder(r)

	_, err := d.Token()
	if err != nil {
		return "", err
	}

	token, err := d.Token()
	if err != nil {
		return "", err
	}

	dataFormat := "jsonl"
	if token == "items" {
		dataFormat = "json"
	}

	return dataFormat, nil
}
