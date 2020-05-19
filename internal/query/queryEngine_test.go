package query

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/denisa/clq/internal/changelog"
	"github.com/stretchr/testify/require"
)

func TestUnsupportedChangelogQuery(t *testing.T) {
	_, err := NewQueryEngine("publication_date")
	require.Error(t, err)
}

func TestUnsupportedChangelogTitleQuery(t *testing.T) {
	_, err := NewQueryEngine("title.size")
	require.Error(t, err)
}

func TestUnsupportedReleaseQuery(t *testing.T) {
	_, err := NewQueryEngine("releases[3].publication_date")
	require.Error(t, err)
}

func TestUnsupportedReleaseIndexQuery(t *testing.T) {
	_, err := NewQueryEngine("releases[three].date")
	require.Error(t, err)
}

func TestQueryTitle(t *testing.T) {
	require := require.New(t)

	result, err := apply("title", []changelog.Heading{
		newHeading(changelog.TitleHeading, "changelog"),
	})
	require.NoError(err)
	require.Equal("changelog", result)
}

func TestQueryTitleAgainstRelease(t *testing.T) {
	require := require.New(t)

	result, err := apply("title", []changelog.Heading{
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
	})
	require.NoError(err)
	require.Empty(result)
}

func TestQueryReleaseAgainstChangelog(t *testing.T) {
	require := require.New(t)

	result, err := apply("releases[1].version", []changelog.Heading{
		newHeading(changelog.TitleHeading, "changelog"),
		newHeading(changelog.TitleHeading, "there should not be another level 1 title"),
	})
	require.NoError(err)
	require.Empty(result)
}

func TestQuerySecondReleaseVersion(t *testing.T) {
	require := require.New(t)

	result, err := apply("releases[1].version", []changelog.Heading{
		newHeading(changelog.TitleHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-15 Cabrel"),
	})
	require.NoError(err)
	require.Equal("1.2.2", result)
}

func TestQuerySecondReleaseDate(t *testing.T) {
	require := require.New(t)

	result, err := apply("releases[0].date", []changelog.Heading{
		newHeading(changelog.TitleHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-15 Cabrel"),
	})
	require.NoError(err)
	require.Equal("2020-05-16", result)
}

func TestQuerySecondReleaseLabel(t *testing.T) {
	require := require.New(t)

	result, err := apply("releases[1].label", []changelog.Heading{
		newHeading(changelog.TitleHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-15 Cabrel"),
	})
	require.NoError(err)
	require.Equal("Cabrel", result)
}

func TestQuerySecondReleaseStatusReleased(t *testing.T) {
	require := require.New(t)

	result, err := apply("releases[1].status", []changelog.Heading{
		newHeading(changelog.TitleHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-15 Cabrel"),
	})
	require.NoError(err)
	require.Equal("released", result)
}

func TestQuerySecondReleaseStatusUnreleased(t *testing.T) {
	require := require.New(t)

	result, err := apply("releases[0].status", []changelog.Heading{
		newHeading(changelog.TitleHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[Unreleased]"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-15 Cabrel"),
	})
	require.NoError(err)
	require.Equal("unreleased", result)
}

func TestQuerySecondReleaseStatusYanked(t *testing.T) {
	require := require.New(t)

	result, err := apply("releases[1].status", []changelog.Heading{
		newHeading(changelog.TitleHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ReleaseHeading, "1.2.2 - 2020-05-15 [YANKED]"),
	})
	require.NoError(err)
	require.Equal("yanked", result)
}

func TestQuerySecondRelease(t *testing.T) {
	require := require.New(t)

	result, err := apply("releases[1]", []changelog.Heading{
		newHeading(changelog.TitleHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-15 Cabrel"),
	})
	require.NoError(err)
	require.JSONEq("{\"version\":\"1.2.2\", \"date\":\"2020-05-15\"}", result)
}

func newHeading(level changelog.HeadingKind, text string) changelog.Heading {
	h, _ := changelog.NewHeading(level, text)
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
