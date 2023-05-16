package crossrefindexer

import (
	"testing"
	"time"

	"github.com/matryer/is"
)

var (
	given1, given2, given3    = "given1", "given2", "given3"
	family1, family2, family3 = "f1", "f2", "f3"
	seq1, seq2, seq3          = "first", "second", "third"
	shortContainerTitle       = []string{"Short Container Title 1", "Short Container Title 2"}
)

// var publishedOnline DateParts{DateParts:[][]int{{2013, 4, 8}}}
var (
	publishedOnline = DateParts{}
	publishedPrint  = DateParts{}
)

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

func generateCrossref(modifiers ...func(*Crossref)) *Crossref {
	ref := Crossref{
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

	for _, modifier := range modifiers {
		modifier(&ref)
	}

	return &ref
}

func generateOutput(modifiers ...func(*SimplifiedPublication)) SimplifiedPublication {
	pub := SimplifiedPublication{
		Title:               []string{"title 1", "title 2"},
		DOI:                 "DOI",
		First_Page:          "200",
		Journal:             []string{"Container Title 1", "Container Title 2"},
		Abbreviated_Journal: []string{"Short Container Title 1", "Short Container Title 2"},
		Volume:              "Volume",
		Issue:               "Issue",
		Year:                2006,
		Bibliographic:       "f1 f2 f3 title 1 Container Title 1 Container Title 2 Short Container Title 1 Short Container Title 2 Volume Issue 200 2006",
	}

	for _, modifier := range modifiers {
		modifier(&pub)
	}

	return pub
}

func Test_ToSimplifiedPublication(t *testing.T) {
	tests := []struct {
		name    string
		input   *Crossref
		want    SimplifiedPublication
		wantErr bool
	}{
		{
			name:    "happy path",
			input:   generateCrossref(),
			want:    generateOutput(),
			wantErr: false,
		},
		{
			name:  "No date",
			input: generateCrossref(func(cr *Crossref) { cr.Issued = DateParts{} }),
			want: generateOutput(func(sp *SimplifiedPublication) {
				sp.Year = 0
				sp.Bibliographic = "f1 f2 f3 title 1 Container Title 1 Container Title 2 Short Container Title 1 Short Container Title 2 Volume Issue 200 0"
			}),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			is.Equal(ToSimplifiedPublication(tt.input), tt.want)
		})
	}
}
