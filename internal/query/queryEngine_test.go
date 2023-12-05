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

func TestParseNameFormatError(t *testing.T) {
	{
		_, _, _, _, err := parseName("changes[")
		require.Error(t, err)
	}
	{
		_, _, _, _, err := parseName("changes]")
		require.Error(t, err)
	}
	{
		_, _, _, _, err := parseName("changes][")
		require.Error(t, err)
	}
	{
		_, _, _, _, err := parseName("changes/")
		require.Error(t, err)
	}
}

func TestParseNameNotRecursiveList(t *testing.T) {
	assertions := require.New(t)
	name, selector, isList, isRecursive, err := parseName("title")
	assertions.NoError(err)
	assertions.Equal("title", name)
	assertions.Equal("", selector)
	assertions.False(isList)
	assertions.False(isRecursive)
}

func TestParseNameListNoSelector(t *testing.T) {
	assertions := require.New(t)
	name, selector, isList, isRecursive, err := parseName("changes[]")
	assertions.NoError(err)
	assertions.Equal("changes", name)
	assertions.Equal("", selector)
	assertions.True(isList)
	assertions.False(isRecursive)
}

func TestParseNameListWithSelector(t *testing.T) {
	assertions := require.New(t)
	name, selector, isList, isRecursive, err := parseName("changes[2]")
	assertions.NoError(err)
	assertions.Equal("changes", name)
	assertions.Equal("2", selector)
	assertions.True(isList)
	assertions.False(isRecursive)
}

func TestParseNameRecursiveListNoSelector(t *testing.T) {
	assertions := require.New(t)
	name, selector, isList, isRecursive, err := parseName("changes[]/")
	assertions.NoError(err)
	assertions.Equal("changes", name)
	assertions.Equal("", selector)
	assertions.True(isList)
	assertions.True(isRecursive)
}

func TestElementIsCollectionInScalarContext(t *testing.T) {
	require.Error(t, elementIsCollection("releases", false))
}

func TestElementIsCollectionNoError(t *testing.T) {
	require.NoError(t, elementIsCollection("releases", true))
}

func TestElementIsFinalNoError(t *testing.T) {
	require.NoError(t, elementIsFinal("title", false, []string{}))
}

func TestElementIsFinalListInScalarContext(t *testing.T) {
	require.Error(t, elementIsFinal("title", true, []string{}))
}

func TestElementIsFinalMoreelementsInScalarContext(t *testing.T) {
	require.Error(t, elementIsFinal("title", false, []string{"kind"}))
}

func newQueryEngine(query string, formatName string) (*QueryEngine, error) {
	outputFormat, err := output.NewFormat(formatName)
	if err != nil {
		return nil, err
	}
	qe, err := NewQueryEngine(query, outputFormat)
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
