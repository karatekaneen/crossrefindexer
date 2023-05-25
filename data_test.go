package crossrefindexer

import (
	"io"
	"log"
	"strings"
	"sync"
	"testing"

	"github.com/matryer/is"
	"go.uber.org/zap"
)

func Test_Load(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		dir         string
		format      Format
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
			format:      FormatJSON,
			compression: "none",
			wantData:    "foxtrot unicorn",
			want: []DataContainer{
				{
					Format:      FormatJSON,
					Compression: "none",
					Data:        io.NopCloser(strings.NewReader("foxtrot unicorn")),
				},
			},
		},
		{
			name:        "file happy path",
			path:        "test.json",
			format:      FormatJSON,
			compression: "none",
			want: []DataContainer{
				{
					Format:      FormatJSON,
					Compression: "none",
					Path:        "test.json",
				},
			},
		},
		{
			name:   "dir happy path with detected compression",
			dir:    "testdata/2021",
			format: FormatJSON,
			want: []DataContainer{
				{
					Format:      FormatJSON,
					Compression: "none",
					Path:        "testdata/2021/0.json",
				},
				{
					Format:      FormatJSON,
					Compression: "gzip",
					Path:        "testdata/2021/0.json.gz",
				},
				{
					Format:      FormatJSON,
					Compression: "gzip",
					Path:        "testdata/2021/1.json.gz",
				},
				{
					Format:      FormatJSON,
					Compression: "gzip",
					Path:        "testdata/2021/2.json.gz",
				},
			},
		},
		{
			name:        "dir happy path with explicit compression",
			dir:         "testdata/2021",
			format:      FormatJSON,
			compression: "none",
			want: []DataContainer{
				{
					Format:      FormatJSON,
					Compression: "none",
					Path:        "testdata/2021/0.json",
				},
				{
					Format:      FormatJSON,
					Compression: "none",
					Path:        "testdata/2021/0.json.gz",
				},
				{
					Format:      FormatJSON,
					Compression: "none",
					Path:        "testdata/2021/1.json.gz",
				},
				{
					Format:      FormatJSON,
					Compression: "none",
					Path:        "testdata/2021/2.json.gz",
				},
			},
		},
		{
			name:   "dir detect both format and compression",
			dir:    "testdata/2021",
			format: FormatUnknown,
			want: []DataContainer{
				{
					Format:      FormatJSON,
					Compression: "none",
					Path:        "testdata/2021/0.json",
				},
				{
					Format:      FormatJSON,
					Compression: "gzip",
					Path:        "testdata/2021/0.json.gz",
				},
				{
					Format:      FormatJSON,
					Compression: "gzip",
					Path:        "testdata/2021/1.json.gz",
				},
				{
					Format:      FormatJSON,
					Compression: "gzip",
					Path:        "testdata/2021/2.json.gz",
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
				is.Equal(item.Format, got[i].Format)
				is.Equal(item.Compression, got[i].Compression)
				is.Equal(item.Path, got[i].Path)
			}

			if tt.data != nil {
				gotData, err := io.ReadAll(got[0].Data)
				is.NoErr(err)
				is.Equal(string(gotData), tt.wantData)

			}
		})
	}
}

func Test_ParseData(t *testing.T) {
	tests := []struct {
		name              string
		input             DataContainer
		wantNumberOfItems int
		wantErr           bool
	}{
		{
			name: "happy path - NDJSON",
			input: DataContainer{
				Format:      "ndjson",
				Compression: "none",
				Path:        "testdata/gap/D1000000.json",
			},
			wantNumberOfItems: 1000,
		},
		{
			name: "happy path - JSON 2022",
			input: DataContainer{
				Format:      "json",
				Compression: "none",
				Path:        "testdata/2022/0.json",
			},
			wantNumberOfItems: 5000,
			wantErr:           false,
		},
		{
			name: "happy path - JSON gzip 2023",
			input: DataContainer{
				Format:      "json",
				Compression: "gzip",
				Path:        "testdata/2023/1.json.gz",
			},
			wantNumberOfItems: 5000,
			wantErr:           false,
		},
		{
			name: "happy path - JSON gzip 2022",
			input: DataContainer{
				Format:      "json",
				Compression: "gzip",
				Path:        "testdata/2022/0.json.gz",
			},
			wantNumberOfItems: 5000,
			wantErr:           false,
		},
		{
			name: "happy path - JSON 2021",
			input: DataContainer{
				Format:      "json",
				Compression: "none",
				Path:        "testdata/2021/0.json",
			},
			wantNumberOfItems: 3000,
			wantErr:           false,
		},
		{
			name: "Wrong type of compression",
			input: DataContainer{
				Format:      "json",
				Compression: "none",
				Path:        "testdata/2022/0.json.gz",
			},
			wantErr: true,
		},
		{
			name: "invalid json",
			input: DataContainer{
				Format:      "json",
				Compression: "none",
				Path:        "testdata/2021/0.json",
				Data:        io.NopCloser(strings.NewReader(`{"someField" = "someValue"}`)),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			ch := make(chan Crossref)

			parsed := 0
			var wg sync.WaitGroup

			wg.Add(1)
			go func() {
				for {
					pub, open := <-ch
					if !open {
						wg.Done()
						return
					}
					parsed++

					// HACK: If we fail the test here it deadlocks because of the posting to
					// the channel within ParseData below. Should be improved
					if pub.Doi == "" {
						log.Fatalf("Invalid data: %+v", pub)
					}
				}
			}()

			err := ParseData(tt.input, ch)
			if tt.wantErr {
				is.True(err != nil)
				return
			}
			close(ch)

			is.NoErr(err)
			wg.Wait()
			is.Equal(parsed, tt.wantNumberOfItems)
		})
	}
}
