package crossrefindexer

import (
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/matryer/is"
)

// func Test_JsonParser(t *testing.T) {
// 	JsonParser()
// }

func Test_JsonParser(t *testing.T) {
	fileOrKill := func(path string) *os.File {
		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}

		return f
	}

	tests := []struct {
		name        string
		data        io.ReadCloser
		fileFormat  string
		itemsInFile int
		numOfLines  int
		wantErr     bool
	}{
		{
			name:        "happy path - JSONL",
			data:        fileOrKill("testdata/gap/D1000000.json"),
			fileFormat:  "jsonl",
			itemsInFile: 1000,
			wantErr:     false,
		},
		{
			name:        "happy path - JSON 2022",
			data:        fileOrKill("testdata/2022/0.json"),
			fileFormat:  "json",
			itemsInFile: 5000,
			wantErr:     false,
		},
		{
			name:        "happy path - JSON 2021",
			data:        fileOrKill("testdata/2021/0.json"),
			fileFormat:  "json",
			itemsInFile: 3000,
			wantErr:     false,
		},
		{
			name:    "invalid json",
			data:    io.NopCloser(strings.NewReader(`{"someField" = "someValue"}`)),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			defer tt.data.Close()

			ch := make(chan Crossref)

			parsed := 0
			var wg sync.WaitGroup

			wg.Add(1)
			go func() {
				for {
					_, open := <-ch
					if !open {
						wg.Done()
						return
					}
					parsed++
				}
			}()

			err := JsonParser(tt.data, ch, tt.fileFormat)
			if tt.wantErr {
				is.True(err != nil)
				return
			}
			close(ch)

			is.NoErr(err)
			wg.Wait()
			is.Equal(parsed, tt.itemsInFile)
		})
	}
}
