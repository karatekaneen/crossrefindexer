package crossrefindexer

import (
	"encoding/json"
	"os"
	"sync"
	"testing"

	"github.com/matryer/is"
)

func Test_Load(t *testing.T) {
	tests := []struct {
		name             string
		path             string
		wantErr          bool
		wantPublications int
	}{
		{
			name:             "happy path - file",
			path:             "testdata/gap/D1000000.json.gz",
			wantErr:          false,
			wantPublications: 1000,
		},
		{
			name:             "happy path - dir",
			path:             "testdata/gap/",
			wantPublications: 3000,
			wantErr:          false,
		},
		{
			name:    "invalid path",
			path:    "testdata/app/",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			ch := make(chan Crossref)

			var wg sync.WaitGroup
			wg.Add(1)
			gotPublications := 0
			go func() {
				for {
					_, stillOpen := <-ch
					if !stillOpen {
						wg.Done()
						return
					}
					gotPublications++
				}
			}()

			err := Load(tt.path, ch)
			if tt.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)
			wg.Wait()
			is.Equal(gotPublications, tt.wantPublications)
		})
	}
}

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

func Test_ClassifyDataFormat(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		jsonType string
		wantErr  bool
	}{
		{
			name:     "happy path",
			file:     "testdata/2022/0.json",
			jsonType: "json",
			wantErr:  false,
		},
		{
			name:    "sad path",
			file:    "testdata/gap/D1000000.json",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			file, err := os.Open(tt.file)
			if err != nil {
				is.NoErr(err)
			}

			defer file.Close()

			jsonTypef, err := ClassifyDataFormat(file)
			is.NoErr(err)

			if tt.wantErr {
				is.True(jsonTypef != tt.jsonType)
				return
			}

			is.True(jsonTypef == tt.jsonType)
		})
	}
}
