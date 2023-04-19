package crossrefindexer

import (
	"github.com/matryer/is"
	"testing"
	"time"
)

func Test_GetSimpleFirstPage(t *testing.T) {
	crossRef_test := CrossRef{
		Page: "389-399",
	}

	tests := []struct {
		name    string
		page    string
		wantErr bool
	}{
		{
			name:    "happy path",
			page:    "389",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			is.Equal(GetSimpleFirstPage(&crossRef_test), tt.page)
			// if tt.wantErr {
			// 	is.True(tt.page != GetSimpleFirstPage(&pub))
			// 	return
			// }

		})
	}

}

func Test_GetSimpleFirstAuthor(t *testing.T) {
	tests := []struct {
		name    string
		author  string
		wantErr bool
	}{
		{
			name:    "happy path",
			author:  "f1",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			is.Equal(GetSimpleFirstAuthor(&crossRef_test), tt.author)

		})
	}
}

func Test_GetSimpleAuthor(t *testing.T) {
	tests := []struct {
		name    string
		author  string
		wantErr bool
	}{
		{
			name:    "happy path",
			author:  "f1 f2 f3",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			is.Equal(GetSimpleAuthor(&crossRef_test), tt.author)

		})
	}
}

func Test_GetSimpleYear(t *testing.T) {
	tests := []struct {
		name    string
		year    int
		wantErr bool
	}{
		{
			name:    "happy path",
			year:    2006,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			is.Equal(GetSimpleYear(&crossRef_test), tt.year)

		})
	}
}

func Test_BuildBibliographicField(t *testing.T) {
	tests := []struct {
		name         string
		bibliography string
		wantErr      bool
	}{
		{
			name:         "happy path",
			bibliography: "f1 f2 f3 title 1 Container Title 1 Container Title 2 Short Container Title 1 Short Container Title 2 Volume Issue 200 2006",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			is.Equal(BuildBibliographicField(&crossRef_test), tt.bibliography)

		})
	}
}

var given1, given2, given3 = "given1", "given2", "given3"
var family1, family2, family3 = "f1", "f2", "f3"
var seq1, seq2, seq3 = "first", "second", "third"
var shortContainerTitle = []string{"Short Container Title 1", "Short Container Title 2"}

// var publishedOnline DateParts{DateParts:[][]int{{2013, 4, 8}}}
var publishedOnline = DateParts{}
var publishedPrint = DateParts{}

var author1 = Author{
	Given:    &given1,
	Family:   &family1,
	Sequence: &seq1,
}
var author2 = Author{
	Given:    &given2,
	Family:   &family2,
	Sequence: &seq2,
}
var author3 = Author{
	Given:    &given3,
	Family:   &family3,
	Sequence: &seq3,
}

var crossRef_test = CrossRef{
	Title:               []string{"title 1", "title 2"},
	Author:              []Author{author1, author2, author3},
	Doi:                 "DOI",
	ContainerTitle:      []string{"Container Title 1", "Container Title 2"},
	ShortContainerTitle: &shortContainerTitle,
	Volume:              "Volume",
	Issue:               "Issue",
	Issued:              DateParts{DateParts: [][]int{{2006, 2, 27}}},
	PublishedOnline:     &publishedOnline,
	PublishedPrint:      &publishedPrint,
	Created: Indexed{
		DateParts: [][]int{{2006, 2, 27}},
		DateTime:  time.Date(2006, time.February, 27, 21, 28, 23, 0, time.UTC),
		Timestamp: 1141075703000,
	},
	Page: "200-300",
}
