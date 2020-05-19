package query

import (
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
