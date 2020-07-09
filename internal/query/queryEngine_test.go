package query

import (
	"testing"

	"github.com/denisa/clq/internal/changelog"
	"github.com/stretchr/testify/require"
)

func TestUnsupportedOutputFormat(t *testing.T) {
	_, err := NewQueryEngine("title", "yaml")
	require.Error(t, err)
}

func TestEmptyQueryAgainstIntroduction(t *testing.T) {
	require := require.New(t)

	result, err := apply("", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
	})
	require.NoError(err)
	require.Equal("", result)
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
	require := require.New(t)
	name, selector, isList, isRecursive, err := parseName("title")
	require.NoError(err)
	require.Equal("title", name)
	require.Equal("", selector)
	require.False(isList)
	require.False(isRecursive)
}

func TestParseNameListNoSelector(t *testing.T) {
	require := require.New(t)
	name, selector, isList, isRecursive, err := parseName("changes[]")
	require.NoError(err)
	require.Equal("changes", name)
	require.Equal("", selector)
	require.True(isList)
	require.False(isRecursive)
}

func TestParseNameListWithSelector(t *testing.T) {
	require := require.New(t)
	name, selector, isList, isRecursive, err := parseName("changes[2]")
	require.NoError(err)
	require.Equal("changes", name)
	require.Equal("2", selector)
	require.True(isList)
	require.False(isRecursive)
}

func TestParseNameRecursiveListNoSelector(t *testing.T) {
	require := require.New(t)
	name, selector, isList, isRecursive, err := parseName("changes[]/")
	require.NoError(err)
	require.Equal("changes", name)
	require.Equal("", selector)
	require.True(isList)
	require.True(isRecursive)
}

func TestElementIsFinalNoerror(t *testing.T) {
	require.NoError(t, elementIsFinal("title", false, []string{}))
}

func TestElementIsFinalListInScalarContext(t *testing.T) {
	require.Error(t, elementIsFinal("title", true, []string{}))
}

func TestElementIsFinalMoreelementsInScalarContext(t *testing.T) {
	require.Error(t, elementIsFinal("title", false, []string{"kind"}))
}

func newHeading(kind changelog.HeadingKind, text string) changelog.Heading {
	h, err := changelog.NewHeading(kind, text)
	if err != nil {
		panic(err)
	}
	return h
}

func apply(query string, headings []changelog.Heading) (string, error) {
	qe, err := NewQueryEngine(query, "json")
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
