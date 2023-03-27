package main

import (
	"github.com/matryer/is"
	"io"
	"log"
	"os"
	"testing"
)

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
			name:     "sad path",
			file:     "testdata/gap/D1000000.json",
			jsonType: "json",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			file, err := os.Open(tt.file)

			if err != nil {
				log.Fatal(err)
			}

			defer file.Close()

			jsonTypef := ClassifyDataFormat(file)
			po, _ := file.Seek(0, io.SeekCurrent)

			if tt.wantErr {
				is.True(jsonTypef != tt.jsonType)
				is.Equal(po, int64(0))
				return
			}

			is.True(jsonTypef == tt.jsonType)
			is.Equal(po, int64(0))
		})
	}
}
