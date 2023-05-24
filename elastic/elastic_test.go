package elastic

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/karatekaneen/crossrefindexer/elastictest"
	"github.com/matryer/is"
	"go.uber.org/zap"
)

func TestDeleteIndex(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		transport *elastictest.ElasticTransport
		want      bool
		wantErr   bool
	}{
		{
			name:  "happy path",
			input: "crossref",
			transport: elastictest.New(
				elastictest.WithResponse(elastictest.CaseDeleteIndexOk),
				elastictest.WithValidation(func(r *http.Request) error {
					if r.URL.Path != "/crossref" || r.Method != "DELETE" {
						return fmt.Errorf(
							"URL %q or Method %q not matching expected",
							r.URL,
							r.Method,
						)
					}
					return nil
				}),
			),
		},
		{
			name:  "index not found",
			input: "crossref",
			transport: elastictest.New(
				elastictest.WithResponse(elastictest.CaseDeleteIndexNotFound),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			idx, err := New(Config{}, zap.NewNop().Sugar(), WithTransport(tt.transport))
			is.NoErr(err)

			err = idx.DeleteIndex(context.Background(), tt.input)
			if tt.wantErr {
				is.True(err != nil)
				return
			}

			is.NoErr(err)
		})
	}
}

func TestCreateIndex(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		transport *elastictest.ElasticTransport
		want      bool
		wantErr   bool
	}{
		{
			name:  "happy path",
			input: "crossref",
			transport: elastictest.New(
				elastictest.WithResponse(elastictest.CaseCreateIndexOk),
				elastictest.WithValidation(func(r *http.Request) error {
					if r.URL.Path != "/crossref" || r.Method != http.MethodPut {
						return fmt.Errorf(
							"URL %q or Method %q not matching expected",
							r.URL,
							r.Method,
						)
					}
					return nil
				}),
			),
		},
		{
			name:  "happy path",
			input: "crossref",
			transport: elastictest.New(
				elastictest.WithResponse(elastictest.CaseCreateIndexConflict),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)

			idx, err := New(Config{}, zap.NewNop().Sugar(), WithTransport(tt.transport))
			is.NoErr(err)

			err = idx.CreateIndex(context.Background(), tt.input, DefaultSettings())
			if tt.wantErr {
				is.True(err != nil)
				return
			}

			is.NoErr(err)
		})
	}
}
