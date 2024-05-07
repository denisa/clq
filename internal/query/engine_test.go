package query

import (
	"testing"

	"github.com/denisa/clq/internal/changelog"
	"github.com/denisa/clq/internal/output"
	"github.com/stretchr/testify/require"
)

func TestEmptyQueryAgainstIntroduction(t *testing.T) {
	assertions := require.New(t)

	result, err := apply("", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
	})
	assertions.NoError(err)
	assertions.Equal("", result)
}

func newQueryEngine(query string, formatName string) (*Engine, error) {
	outputFormat, err := output.NewFormat(formatName)
	if err != nil {
		return nil, err
	}
	qe, err := NewEngine(query, outputFormat)
	if err != nil {
		return nil, err
	}
	return qe, nil
}

func newHeading(kind changelog.HeadingKind, text string) changelog.Heading {
	ck, _ := changelog.NewChangeKind("")
	hf := changelog.NewHeadingFactory(ck)
	h, err := hf.NewHeading(kind, text)
	if err != nil {
		panic(err)
	}
	return h
}

func apply(query string, headings []changelog.Heading) (string, error) {
	qe, err := newQueryEngine(query, "json")
	if err != nil {
		return "", err
	}

	var stack []changelog.Heading
	for _, h := range headings {
		i := len(stack) - 1
		for ; i >= 0 && stack[i].Kind() >= h.Kind(); i-- {
			qe.Exit(stack[i])
		}
		stack = append(stack[:i+1], h)
		qe.Enter(h)
	}
	for i := len(stack) - 1; i >= 0; i-- {
		qe.Exit(stack[i])
	}
	return qe.Result(), nil
}
