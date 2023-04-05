package crossrefindexer

import (
	"testing"

	"sync"

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
			ch := make(chan CrossRef)

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