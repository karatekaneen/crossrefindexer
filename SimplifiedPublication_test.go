package crossrefindexer

import (
	"github.com/matryer/is"
	"testing"
	"time"
)

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

var simplifiedPublication1 = SimplifiedPublication{
	title:               []string{"title 1", "title 2"},
	DOI:                 "DOI",
	first_page:          "200",
	journal:             []string{"Container Title 1", "Container Title 2"},
	abbreviated_journal: []string{"Short Container Title 1", "Short Container Title 2"},
	volume:              "Volume",
	issue:               "Issue",
	year:                2006,
	Bibliographic:       "f1 f2 f3 title 1 Container Title 1 Container Title 2 Short Container Title 1 Short Container Title 2 Volume Issue 200 2006",
}

func Test_ToSimplifiedPublication(t *testing.T) {
	tests := []struct {
		name    string
		simpPub SimplifiedPublication
		wantErr bool
	}{
		{
			name:    "happy path",
			simpPub: simplifiedPublication1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			is.Equal(ToSimplifiedPublication(&crossRef_test), tt.simpPub)

		})
	}
}
