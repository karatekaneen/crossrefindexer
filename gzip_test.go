package main

import (
	"encoding/json"
	"github.com/matryer/is"
	"os"
	"testing"
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
		{
			name:    "sad path",
			file:    "testdata/gap/D1000000.json.gz",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			file, err := os.Open(tt.file)

			if tt.wantErr {
				is.True(err != nil)
				return
			}

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
			err = j.Decode(&elm)
			is.NoErr(err) // Decoded frist Element
		})
	}
}
