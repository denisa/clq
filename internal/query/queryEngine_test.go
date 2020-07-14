package query

import (
	"github.com/denisa/clq/internal/changelog"
)

func newHeading(kind changelog.HeadingKind, text string) changelog.Heading {
	h, _ := changelog.NewHeading(kind, text)
	return h
}

func apply(query string, headings []changelog.Heading) (string, error) {
	qe, err := NewQueryEngine(query)
	if err != nil {
		return "", err
	}

	for _, h := range headings {
		qe.Enter(h)
	}
	for i := len(headings); i > 0; {
		i--
		qe.Exit(headings[i])
	}
	return qe.Result(), nil
}
