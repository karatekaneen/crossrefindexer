package crossrefindexer

import (
	"github.com/matryer/is"
	"testing"
)

func Test_Load(t *testing.T) {

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "happy path - file",
			path:    "testdata/gap/D1000000.json.gz",
			wantErr: false,
		},
		{
			name:    "happy path - dir",
			path:    "testdata/gap/",
			wantErr: false,
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

			err := Load(tt.path)
			if tt.wantErr {
				is.True(err != nil)
				return
			}

			is.NoErr(err)
		})
	}
}
