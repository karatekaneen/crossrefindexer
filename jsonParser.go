package crossrefindexer

import (
	"encoding/json"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func JsonParser(r io.Reader, ch chan CrossRef, format string) error {
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
	Abstract            *string       `json:"abstract"` // Gap
	Author              []Author      `json:"author"`
	ContainerTitle      []string      `json:"container-title"`
	ContentDomain       ContentDomain `json:"content-domain"`
	Created             Indexed       `json:"created"`
	Deposited           Indexed       `json:"deposited"`
	Doi                 string        `json:"DOI"`
	Indexed             Indexed       `json:"indexed"`
	IsReferencedByCount int           `json:"is-referenced-by-count"`
	Issn                []string      `json:"ISSN"`
	IssnType            []IssnType    `json:"issn-type"`
	Issue               string        `json:"issue"`
	Issued              DateParts     `json:"issued"`
	JournalIssue        *JournalIssue `json:"journal-issue"` // Gap
	Language            *string       `json:"language"`      //Gap
	Link                *[]Link       `json:"link"`          // Gap
	Member              string        `json:"member"`
	OriginalTitle       *[]any        `json:"original-title"` // 2021
	Page                string        `json:"page"`
	Prefix              string        `json:"prefix"`
	Published           *DateParts    `json:"published"`        // Gap
	PublishedOnline     *DateParts    `json:"published-online"` // Gap
	PublishedOther      *DateParts    `json:"published-other"`  // Gap
	PublishedPrint      *DateParts    `json:"published-print"`
	Publisher           string        `json:"publisher"`
	Reference           *[]Reference  `json:"reference"` // Gap
	ReferenceCount      int           `json:"reference-count"`
	ReferencesCount     int           `json:"references-count"`
	Relation            *Relation     `json:"relation"` // 2021
	Resource            Resource      `json:"resource"`
	Score               float64       `json:"score"`
	ShortContainerTitle *[]string     `json:"short-container-title"` // 2021
	ShortTitle          *[]any        `json:"short-title"`           // 2021
	Source              string        `json:"source"`
	Subject             []string      `json:"subject"`
	Subtitle            *[]any        `json:"subtitle"` // 2021
	Title               []string      `json:"title"`
	Type                string        `json:"type"`
	URL                 string        `json:"URL"`
	UpdatePolicy        *string       `json:"update-policy"` // Gap
	Volume              string        `json:"volume"`
	License             []License     `json:"license"`
	AlternativeID       []string      `json:"alternative-id"`
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
type Affiliation struct {
	Name string `json:"name"`
}
type Author struct {
	Given       *string        `json:"given"`
	Family      *string        `json:"family"`
	Sequence    *string        `json:"sequence"`
	Affiliation *[]Affiliation `json:"affiliation"`
}
type DateParts struct {
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
type JournalIssue struct {
	Issue           *string    `json:"issue"`
	PublishedOnline *DateParts `json:"published-online"`
	PublishedPrint  *DateParts `json:"published-print"`
}
type IssnType struct {
	Value string `json:"value"`
	Type  string `json:"type"`
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

//
func (pub *CrossRef) GetSimpleTitle() []string {
	simpleTitle := pub.Title
	if len(simpleTitle) == 0 {
		return []string{""}
	}
	for pos, title := range simpleTitle {
		title = strings.Replace(title, "\n", " ", -1)
		title = strings.Replace(title, "( )+", " ", -1)
		simpleTitle[pos] = strings.TrimSpace(title)
	}
	return simpleTitle
}

func (pub *CrossRef) GetSimpleDOI() string {
	return pub.Doi
}

func (pub *CrossRef) GetSimpleAuthor() string {
	if len(pub.Author) == 0 {
		return ""
	}
	var simpleAuthor strings.Builder
	for _, auth := range pub.Author {
		simpleAuthor.WriteString(*auth.Family + " ")
	}
	return strings.TrimSpace(simpleAuthor.String())
}

func (pub *CrossRef) GetSimpleFirstAuthor() string {
	if len(pub.Author) == 0 {
		return ""
	}
	for _, auth := range pub.Author {
		if *auth.Sequence == "first" && *auth.Family != "" {
			return *auth.Family
		}
	}
	// not sequence information apparently, so as fallback we use the first
	// author in the author list
	return *pub.Author[0].Family
}

func (pub *CrossRef) GetSimpleFirstPage() string {
	sp := regexp.MustCompile(`,|-|\s`)
	pagePieces := sp.Split(pub.Page, -1)
	return pagePieces[0]
}

func (pub *CrossRef) GetSimpleJournal() string {
	if len(pub.ContainerTitle) == 0 {
		return ""
	}
	return pub.ContainerTitle[0]
}
func (pub *CrossRef) GetSimpleAbbreviatedJournal() string {
	if len(*pub.ShortContainerTitle) == 0 {
		return ""
	}
	return (*pub.ShortContainerTitle)[0]
}

func (pub *CrossRef) GetSimpleVolume() string {
	return pub.Volume
}

func (pub *CrossRef) GetSimpleIssue() string {
	return pub.Issue
}

// year is a date part (first one) in issued or created or published-online (we follow this order)
func (pub *CrossRef) GetSimpleYear() int {
	var year int
	switch {
	case pub.Issued.DateParts != nil:
		year = pub.Issued.DateParts[0][0]
	case pub.PublishedOnline != nil:
		year = pub.PublishedOnline.DateParts[0][0]
	case pub.PublishedPrint != nil:
		year = pub.PublishedPrint.DateParts[0][0]
	default:
		// this is deposit date, normally we will never use it, but it will ensure
		// that we always have a date as conservative fallback
		year = pub.Created.DateParts[0][0]
	}
	return year
}

func BuildBibliographicField(pub *CrossRef) string {
	var res strings.Builder
	if pub.GetSimpleAuthor() != "" {
		res.WriteString(pub.GetSimpleAuthor() + " ")
	} else {
		res.WriteString(pub.GetSimpleFirstAuthor() + " ")
	}
	res.WriteString(pub.GetSimpleTitle()[0] + " ")
	res.WriteString(pub.GetSimpleJournal() + " ")
	res.WriteString(pub.GetSimpleAbbreviatedJournal() + " ")
	res.WriteString(pub.GetSimpleVolume() + " ")
	res.WriteString(pub.GetSimpleIssue() + " ")
	res.WriteString(pub.GetSimpleFirstPage() + " ")

	year := pub.GetSimpleYear()
	res.WriteString(strconv.FormatInt(int64(year), 10) + " ")

	return strings.TrimSpace(res.String())
}
