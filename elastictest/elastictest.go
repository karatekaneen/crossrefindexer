package elastictest

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type ElasticTransport struct {
	validation func(*http.Request) error // To validate the request being made
	response   io.ReadCloser             // The data to send back
	code       int                       // Statuscode to return
}

type (
	Option   func(*ElasticTransport) // To customize the struct
	TestCase struct {
		path       string // Path to the file containing the result
		statusCode int    // The status code to return
	}
)

var (
	CaseCreateIndexOk = TestCase{
		path:       "testdata/elastic/create_index_ok.json",
		statusCode: 200,
	}
	CaseCreateIndexConflict = TestCase{
		path:       "testdata/elastic/create_index_conflict.json",
		statusCode: 400,
	}
	CaseDeleteIndexOk = TestCase{
		path:       "testdata/elastic/delete_index_ok.json",
		statusCode: 200,
	}
	CaseDeleteIndexNotFound = TestCase{
		path:       "testdata/elastic/delete_index_notfound.json",
		statusCode: 404,
	}
)

func WithValidation(val func(*http.Request) error) Option {
	return func(et *ElasticTransport) { et.validation = val }
}

// WithResponse reads the testfile into the response.
// Will crash if file reading goes wrong.
func WithResponse(tCase TestCase) Option {
	root, err := getProjectRoot()
	if err != nil {
		log.Fatal(err) // In test code we crash everything!
	}

	f, err := os.Open(filepath.Join(root, tCase.path))
	if err != nil {
		log.Fatal(err)
	}

	return func(et *ElasticTransport) {
		et.response = f
		et.code = tCase.statusCode
	}
}

func New(options ...Option) *ElasticTransport {
	e := &ElasticTransport{
		validation: func(r *http.Request) error { return nil },
		response:   io.NopCloser(strings.NewReader(`{}`)),
		code:       http.StatusOK,
	}

	for _, option := range options {
		option(e)
	}

	return e
}

// This is a naive solution to be able to open files simply
// when you're running tests from different working directories.
//
// * Will probably crash in some environments *
func getProjectRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "could not get project root path")
	}

	projectName := "crossrefindexer"
	parentFolder, _, _ := strings.Cut(wd, projectName)
	return filepath.Join(parentFolder, projectName), nil
}

func (et *ElasticTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if r.URL.Path == "/" && r.Method == "GET" {
		return response(http.StatusOK, io.NopCloser(strings.NewReader(""))), nil
	}
	if err := et.validation(r); err != nil {
		return nil, err
	}

	// * The header is needed so that the Elastic client won't shit itself
	return response(et.code, et.response), nil
}

func response(status int, body io.ReadCloser) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       body,
		Header:     http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
	}
}
