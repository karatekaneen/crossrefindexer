package crossrefindexer

import (
	"encoding/json"
	"github.com/matryer/is"
	"os"
	"testing"
)

// func Test_JsonParser(t *testing.T) {
// 	JsonParser()
// }

func Test_JsonParser(t *testing.T) {

	tests := []struct {
		name       string
		file       string
		numOfElem  int
		numOfLines int
		wantErr    bool
	}{
		{
			name:    "happy path",
			file:    "testdata/gap/D1000000.json",
			wantErr: false,
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
			jsonDecoder := json.NewDecoder(file)
			jsonType, err := ClassifyDataFormat(file)

			// skip the items array
			if jsonType == "json" {
				jsonDecoder.Token()
				jsonDecoder.Token()
				jsonDecoder.Token()
			}
			for i := 0; i < 2; i++ {
				elm, err := JsonParser(jsonDecoder)
				println(elm, err)
			}
			// var elm map[string]any
			// err = j.Decode(&elm)
			// is.NoErr(err) // Decoded frist Element
		})
	}
}
