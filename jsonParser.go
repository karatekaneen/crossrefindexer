package crossrefindexer

import (
	"encoding/json"
	"io"
	"time"

	"github.com/pkg/errors"
)

func JsonParser(r io.Reader, ch chan CrossRef, format string) error {
	defer close(ch)
	d := json.NewDecoder(r)

	// The json format is quite nested so we need to skip
	// three levels "deep" to reach the data that we want
	if format == "json" {
		d.Token()
		d.Token()
		d.Token()
	}

	elementIndex := 0

	for d.More() {
		var publication CrossRef

		err := d.Decode(&publication)
		if err == io.EOF {
			break
		} else if err != nil {
			return errors.Wrapf(err, "failed on parsing element %d", elementIndex)
		}

		ch <- publication
		elementIndex++
	}
	return nil
}

// Reference and value semantics reflect required and optional value in json
type CrossRef struct {
	Abstract            *string          `json:"abstract"` // Gap
	Author              []Author         `json:"author"`
	ContainerTitle      []string         `json:"container-title"`
	ContentDomain       ContentDomain    `json:"content-domain"`
	Created             Indexed          `json:"created"`
	Deposited           Indexed          `json:"deposited"`
	Doi                 string           `json:"DOI"`
	Indexed             Indexed          `json:"indexed"`
	IsReferencedByCount int              `json:"is-referenced-by-count"`
	Issn                []string         `json:"ISSN"`
	IssnType            []IssnType       `json:"issn-type"`
	Issue               string           `json:"issue"`
	Issued              Issued           `json:"issued"`
	JournalIssue        *JournalIssue    `json:"journal-issue"` // Gap
	Language            *string          `json:"language"`      //Gap
	Link                *[]Link          `json:"link"`          // Gap
	Member              string           `json:"member"`
	OriginalTitle       *[]any           `json:"original-title"` // 2021
	Page                string           `json:"page"`
	Prefix              string           `json:"prefix"`
	Published           *Published       `json:"published"`        // Gap
	PublishedOnline     *PublishedOnline `json:"published-online"` // Gap
	PublishedOther      *PublishedOther  `json:"published-other"`  // Gap
	PublishedPrint      *PublishedPrint  `json:"published-print"`
	Publisher           string           `json:"publisher"`
	Reference           *[]Reference     `json:"reference"` // Gap
	ReferenceCount      int              `json:"reference-count"`
	ReferencesCount     int              `json:"references-count"`
	Relation            *Relation        `json:"relation"` // 2021
	Resource            Resource         `json:"resource"`
	Score               float64          `json:"score"`
	ShortContainerTitle *[]string        `json:"short-container-title"` // 2021
	ShortTitle          *[]any           `json:"short-title"`           // 2021
	Source              string           `json:"source"`
	Subject             []string         `json:"subject"`
	Subtitle            *[]any           `json:"subtitle"` // 2021
	Title               []string         `json:"title"`
	Type                string           `json:"type"`
	URL                 string           `json:"URL"`
	UpdatePolicy        *string          `json:"update-policy"` // Gap
	Volume              string           `json:"volume"`
	License             []License        `json:"license"`
	AlternativeID       []string         `json:"alternative-id"`
}

type Indexed struct {
	DateParts [][]int   `json:"date-parts"`
	DateTime  time.Time `json:"date-time"`
	Timestamp int64     `json:"timestamp"`
}

type ContentDomain struct {
	Domain               []string `json:"domain"`
	CrossmarkRestriction bool     `json:"crossmark-restriction"`
}
type PublishedPrint struct {
	DateParts [][]int `json:"date-parts"`
}

type Affiliation struct {
	Name string `json:"name"`
}
type Author struct {
	Given       *string        `json:"given"`
	Family      *string        `json:"family"`
	Sequence    *string        `json:"sequence"`
	Affiliation *[]Affiliation `json:"affiliation"`
}
type PublishedOnline struct {
	DateParts [][]int `json:"date-parts"`
}
type Reference struct {
	Key           string  `json:"key"`
	VolumeTitle   string  `json:"volume-title,omitempty"`
	Author        string  `json:"author"`
	Year          string  `json:"year"`
	FirstPage     string  `json:"first-page,omitempty"`
	ArticleTitle  string  `json:"article-title,omitempty"`
	DoiAssertedBy string  `json:"doi-asserted-by,omitempty"`
	Doi           string  `json:"DOI,omitempty"`
	Volume        string  `json:"volume,omitempty"`
	JournalTitle  string  `json:"journal-title,omitempty"`
	Issue         string  `json:"issue,omitempty"`
	Unstructured  *string `json:"unstructured,omitempty"`
}
type Link struct {
	URL                 string `json:"URL"`
	ContentType         string `json:"content-type"`
	ContentVersion      string `json:"content-version"`
	IntendedApplication string `json:"intended-application"`
}

type Primary struct {
	URL string `json:"URL"`
}
type Resource struct {
	Primary Primary `json:"primary"`
}
type Issued struct {
	DateParts [][]int `json:"date-parts"`
}
type JournalIssue struct {
	Issue           *string          `json:"issue"`
	PublishedOnline *PublishedOnline `json:"published-online"`
	PublishedPrint  *PublishedPrint  `json:"published-print"`
}
type IssnType struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}
type PublishedOther struct {
	DateParts [][]int `json:"date-parts"`
}
type Published struct {
	DateParts [][]int `json:"date-parts"`
}

type Relation struct {
	Cites []any `json:"cites"`
}

type License struct {
	URL            string  `json:"URL"`
	Start          Indexed `json:"start"`
	DelayInDays    int     `json:"delay-in-days"`
	ContentVersion string  `json:"content-version"`
}