package crossrefindexer

import (
	"fmt"
	"regexp"
	"strings"
)

func pubTitle(pub CrossRef) []string {
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

func firstPage(pub *CrossRef) string {
	sp := regexp.MustCompile(
		`,|-` +
			// This matches any white space character, including spaces, tabs, and newlines.
			`|\s`)
	pagePieces := sp.Split(pub.Page, -1)
	return pagePieces[0]
}

// year is a date part (first one) in issued or created or published-online (we follow this order)
func pubYear(pub *CrossRef) int {
	var year int
	switch {
	case pub.Issued.DateParts != nil:
		year = pub.Issued.DateParts[0][0]
	case pub.PublishedOnline != nil:
		year = pub.PublishedOnline.DateParts[0][0]
	case pub.PublishedPrint != nil:
		year = pub.PublishedPrint.DateParts[0][0]
	case pub.Created.DateParts != nil:
		// this is deposit date, normally we will never use it, but it will ensure
		// that we always have a date as conservative fallback
		year = pub.Created.DateParts[0][0]
	default:
		year = 0
	}
	return year
}

func buildBibliographicField(pub *CrossRef) string {
	author := make([]string, len(pub.Author))
	for _, auth := range pub.Author {
		if *auth.Family == "" {
			continue
		}
		author = append(author, *auth.Family)
	}

	bibliographic := []string{
		strings.TrimSpace(strings.Join(author, " ")),
		pubTitle(*pub)[0],
		strings.Join(pub.ContainerTitle, " "),
		strings.Join(*pub.ShortContainerTitle, " "),
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

func ToSimplifiedPublication(pub *CrossRef) SimplifiedPublication {
	var simpPub SimplifiedPublication
	simpPub.Title = pubTitle(*pub)
	simpPub.DOI = pub.Doi
	simpPub.First_Page = firstPage(pub)
	simpPub.Journal = pub.ContainerTitle
	simpPub.Abbreviated_Journal = *pub.ShortContainerTitle
	simpPub.Volume = pub.Volume
	simpPub.Issue = pub.Issue
	simpPub.Year = pubYear(pub)
	simpPub.Bibliographic = buildBibliographicField(pub)
	return simpPub
}
