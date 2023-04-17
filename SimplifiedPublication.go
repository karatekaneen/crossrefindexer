package crossrefindexer

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

func GetSimpleTitle(pub *CrossRef) []string {
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

func GetSimpleAuthor(pub *CrossRef) string {
	if len(pub.Author) == 0 {
		return ""
	}
	var simpleAuthor strings.Builder
	for _, auth := range pub.Author {
		simpleAuthor.WriteString(*auth.Family + " ")
	}
	return strings.TrimSpace(simpleAuthor.String())
}

func GetSimpleFirstAuthor(pub *CrossRef) string {
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

func SimpleAuthor(pub *CrossRef) string {
	if GetSimpleAuthor(pub) != "" {
		return GetSimpleAuthor(pub)
	} else {
		return GetSimpleFirstAuthor(pub)
	}
}

func GetSimpleFirstPage(pub *CrossRef) string {
	sp := regexp.MustCompile(`,|-|\s`)
	pagePieces := sp.Split(pub.Page, -1)
	return pagePieces[0]
}

// year is a date part (first one) in issued or created or published-online (we follow this order)
func GetSimpleYear(pub *CrossRef) (int, error) {
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
		return 0, errors.New("Faild to reterive a year")
	}
	return year, nil
}

func BuildBibliographicField(pub *CrossRef) string {

	year, _ := GetSimpleYear(pub)

	bibliographic := []string{
		SimpleAuthor(pub),
		GetSimpleTitle(pub)[0],
		strings.Join(pub.ContainerTitle, " "),
		strings.Join(*pub.ShortContainerTitle, " "),
		pub.Volume,
		pub.Issue,
		GetSimpleFirstPage(pub),
		strconv.FormatInt(int64(year), 10),
	}

	return strings.Join(bibliographic, " ")
}

type SimplifiedPublication struct {
	title               []string
	DOI                 string
	first_page          string
	journal             []string
	abbreviated_journal []string
	volume              string
	issue               string
	year                int
	Bibliographic       string
}

func ToSimplifiedPublication(pub *CrossRef) SimplifiedPublication {
	var simplifiedPub SimplifiedPublication
	simplifiedPub.title = GetSimpleTitle(pub)
	simplifiedPub.DOI = pub.Doi
	simplifiedPub.first_page = GetSimpleFirstPage(pub)
	simplifiedPub.journal = pub.ContainerTitle
	simplifiedPub.abbreviated_journal = *pub.ShortContainerTitle
	simplifiedPub.volume = pub.Volume
	simplifiedPub.issue = pub.Issue
	simplifiedPub.year, _ = GetSimpleYear(pub)
	simplifiedPub.Bibliographic = BuildBibliographicField(pub)
	return simplifiedPub
}
