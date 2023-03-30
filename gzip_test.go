package crossrefindexer

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/matryer/is"
)

func Test_GzipReader(t *testing.T) {

	tests := []struct {
		name       string
		file       string
		numOfElem  int
		numOfLines int
		wantErr    bool
	}{
		{
			name:    "happy path",
			file:    "testdata/gap/D1000000.json.gz",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			file, err := os.Open(tt.file)
			is.NoErr(err)

			defer file.Close()

			gzr, err := GzipReader(file)
			is.NoErr(err)

			defer gzr.Close()
			if tt.wantErr {
				is.True(err != nil)
				return
			}
			j := json.NewDecoder(gzr)
			var elm map[string]any
			is.NoErr(j.Decode(&elm)) // Decoded frist Element
			is.True(len(elm) > 0)
		})
	}
}
