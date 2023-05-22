package crossrefindexer

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Reference and value semantics reflect required and optional value in json
type Crossref struct {
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
	Language            *string       `json:"language"`      // Gap
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

func pubTitle(pub Crossref) []string {
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

func firstPage(pub *Crossref) string {
	sp := regexp.MustCompile(
		`,|-` +
			// This matches any white space character, including spaces, tabs, and newlines.
			`|\s`)
	pagePieces := sp.Split(pub.Page, -1)
	return pagePieces[0]
}

// year is a date part (first one) in issued or created or published-online (we follow this order)
func pubYear(pub *Crossref) int {
	var year int
	switch {
	case pub.Issued.DateParts != nil:
		year = extractYear(pub.Issued.DateParts)
	case pub.PublishedOnline != nil:
		year = extractYear(pub.PublishedOnline.DateParts)
	case pub.PublishedPrint != nil:
		year = extractYear(pub.PublishedPrint.DateParts)
	case pub.Created.DateParts != nil:
		// this is deposit date, normally we will never use it, but it will ensure
		// that we always have a date as conservative fallback
		year = extractYear(pub.Created.DateParts)
	default:
		year = 0
	}
	return year
}

// extractYear from dateparts and handle if it is empty
func extractYear(dp [][]int) int {
	if len(dp) < 1 || len(dp[0]) < 1 {
		return 0
	}

	return dp[0][0]
}

func buildBibliographicField(pub *Crossref) string {
	author := make([]string, len(pub.Author))
	for _, auth := range pub.Author {
		if stringFromPointer(auth.Family) == "" {
			continue
		}
		author = append(author, *auth.Family)
	}

	abbreviatedJournal := []string{}
	if pub.ShortContainerTitle != nil {
		abbreviatedJournal = *pub.ShortContainerTitle
	}

	bibliographic := []string{
		strings.TrimSpace(strings.Join(author, " ")),
		pubTitle(*pub)[0],
		strings.Join(pub.ContainerTitle, " "),
		strings.Join(abbreviatedJournal, " "),
		pub.Volume,
		pub.Issue,
		firstPage(pub),
		fmt.Sprint(pubYear(pub)),
	}

	return strings.Join(bibliographic, " ")
}

type SimplifiedPublication struct {
	Title               []string
	DOI                 string
	First_Page          string
	Journal             []string
	Abbreviated_Journal []string
	Volume              string
	Issue               string
	Year                int
	Bibliographic       string
}

func stringFromPointer(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ToSimplifiedPublication(pub *Crossref) SimplifiedPublication {
	abbreviatedJournal := []string{}
	if pub.ShortContainerTitle != nil {
		abbreviatedJournal = *pub.ShortContainerTitle
	}

	var simpPub SimplifiedPublication
	simpPub.Title = pubTitle(*pub)
	simpPub.DOI = pub.Doi
	simpPub.First_Page = firstPage(pub)
	simpPub.Journal = pub.ContainerTitle
	simpPub.Abbreviated_Journal = abbreviatedJournal
	simpPub.Volume = pub.Volume
	simpPub.Issue = pub.Issue
	simpPub.Year = pubYear(pub)
	simpPub.Bibliographic = buildBibliographicField(pub)
	return simpPub
}
