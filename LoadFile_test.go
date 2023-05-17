package crossrefindexer

import (
	"io"
	"strings"
	"testing"

	"github.com/matryer/is"
	"go.uber.org/zap"
)

func Test_Load(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		dir         string
		format      string
		compression string
		data        io.Reader
		want        []DataContainer
		wantData    string
		wantErr     bool
	}{
		{
			name:        "Stdin happy path",
			path:        "-",
			data:        strings.NewReader("foxtrot unicorn"),
			format:      "json",
			compression: "none",
			wantData:    "foxtrot unicorn",
			want: []DataContainer{
				{
					format:      "json",
					compression: "none",
					data:        io.NopCloser(strings.NewReader("foxtrot unicorn")),
				},
			},
		},
		{
			name:        "file happy path",
			path:        "test.json",
			format:      "json",
			compression: "none",
			wantData:    "foxtrot unicorn",
			want: []DataContainer{
				{
					format:      "json",
					compression: "none",
					path:        "test.json",
				},
			},
		},
		{
			name:     "dir happy path with detected compression",
			dir:      "testdata/2021",
			format:   "json",
			wantData: "foxtrot unicorn",
			want: []DataContainer{
				{
					format:      "json",
					compression: "none",
					path:        "testdata/2021/0.json",
				},
				{
					format:      "json",
					compression: "gzip",
					path:        "testdata/2021/0.json.gz",
				},
				{
					format:      "json",
					compression: "gzip",
					path:        "testdata/2021/1.json.gz",
				},
				{
					format:      "json",
					compression: "gzip",
					path:        "testdata/2021/2.json.gz",
				},
			},
		},
		{
			name:        "dir happy path with explicit compression",
			dir:         "testdata/2021",
			format:      "json",
			compression: "none",
			wantData:    "foxtrot unicorn",
			want: []DataContainer{
				{
					format:      "json",
					compression: "none",
					path:        "testdata/2021/0.json",
				},
				{
					format:      "json",
					compression: "none",
					path:        "testdata/2021/0.json.gz",
				},
				{
					format:      "json",
					compression: "none",
					path:        "testdata/2021/1.json.gz",
				},
				{
					format:      "json",
					compression: "none",
					path:        "testdata/2021/2.json.gz",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			got, err := Load(
				zap.NewNop().Sugar(),
				tt.path,
				tt.dir,
				tt.format,
				tt.compression,
				tt.data,
			)
			if tt.wantErr {
				is.True(err != nil)
				return
			}

			is.NoErr(err)
			is.Equal(len(tt.want), len(got))

			for i, item := range tt.want {
				is.Equal(item.format, got[i].format)
				is.Equal(item.compression, got[i].compression)
				is.Equal(item.path, got[i].path)
			}

			if tt.data != nil {
				gotData, err := io.ReadAll(got[0].data)
				is.NoErr(err)
				is.Equal(string(gotData), tt.wantData)

			}
		})
	}
}

// func Test_GzipReader(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		file       string
// 		numOfElem  int
// 		numOfLines int
// 		wantErr    bool
// 	}{
// 		{
// 			name:    "happy path",
// 			file:    "testdata/gap/D1000000.json.gz",
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			is := is.New(t)
//
// 			file, err := os.Open(tt.file)
// 			is.NoErr(err)
//
// 			defer file.Close()
//
// 			gzr, err := GzipReader(file)
// 			is.NoErr(err)
//
// 			defer gzr.Close()
// 			if tt.wantErr {
// 				is.True(err != nil)
// 				return
// 			}
// 			j := json.NewDecoder(gzr)
// 			var elm map[string]any
// 			is.NoErr(j.Decode(&elm)) // Decoded frist Element
// 			is.True(len(elm) > 0)
// 		})
// 	}
// }
//
// func Test_ClassifyDataFormat(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		file     string
// 		jsonType string
// 		wantErr  bool
// 	}{
// 		{
// 			name:     "happy path",
// 			file:     "testdata/2022/0.json",
// 			jsonType: "json",
// 			wantErr:  false,
// 		},
// 		{
// 			name:    "sad path",
// 			file:    "testdata/gap/D1000000.json",
// 			wantErr: true,
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			is := is.New(t)
//
// 			file, err := os.Open(tt.file)
// 			if err != nil {
// 				is.NoErr(err)
// 			}
//
// 			defer file.Close()
//
// 			jsonTypef, err := ClassifyDataFormat(file)
// 			is.NoErr(err)
//
// 			if tt.wantErr {
// 				is.True(jsonTypef != tt.jsonType)
// 				return
// 			}
//
// 			is.True(jsonTypef == tt.jsonType)
// 		})
// 	}
// }
