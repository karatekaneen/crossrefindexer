package main

import (
	"github.com/matryer/is"
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

			if tt.wantErr {
				is.True(err != nil)
				return
			}

			defer file.Close()

			jsonTypef := ClassifyDataFormat(file)

			if tt.wantErr {
				is.True(jsonTypef != tt.jsonType)
				return
			} else {
				is.True(jsonTypef == tt.jsonType)
			}
		})
	}
}
