package query

import (
	"bufio"
	"bytes"

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

	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	for _, h := range headings {
		qe.Apply(w, h)
	}
	w.Flush()
	return buf.String(), nil
}
