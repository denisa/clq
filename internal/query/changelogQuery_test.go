package query

import (
	"testing"

	"github.com/denisa/clq/internal/changelog"
	"github.com/stretchr/testify/require"
)

func TestChangelogQueryUnsupportedAttribue(t *testing.T) {
	_, err := NewQueryEngine("publication_date", "json")
	require.Error(t, err)
}

func TestChangelogQueryUnsupportedTitle(t *testing.T) {
	_, err := NewQueryEngine("title.size", "json")
	require.Error(t, err)
}

func TestChangelogQueryTitle(t *testing.T) {
	require := require.New(t)

	result, err := apply("title", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
	})
	require.NoError(err)
	require.Equal("changelog", result)
}

func TestChangelogQueryTitleAgainstRelease(t *testing.T) {
	require := require.New(t)

	result, err := apply("title", []changelog.Heading{
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
	})
	require.NoError(err)
	require.Empty(result)
}

func TestChangelogQueryUnsupportedEnter(t *testing.T) {
	require := require.New(t)

	query := &changelogQuery{}
	require.False(query.Enter(newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16")))
}

func TestChangelogQueryUnsupportedExit(t *testing.T) {
	require := require.New(t)

	query := &changelogQuery{}
	require.False(query.Exit(newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16")))
}

func TestChangelogQueryCollection(t *testing.T) {
	require := require.New(t)
	{
		query := &changelogQuery{}
		require.False(query.isCollection())
	}
	{
		query := &changelogQuery{projections{collection: true}}
		require.True(query.isCollection())
	}
}
