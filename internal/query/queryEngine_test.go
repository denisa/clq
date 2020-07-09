package query

import (
	"testing"

	"github.com/denisa/clq/internal/changelog"
	"github.com/stretchr/testify/require"
)

func TestEmptyQueryAgainstIntroduction(t *testing.T) {
	require := require.New(t)

	result, err := apply("", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
	})
	require.NoError(err)
	require.Equal("", result)
}

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
