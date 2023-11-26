package query

import (
	"testing"

	"github.com/denisa/clq/internal/changelog"
	"github.com/stretchr/testify/require"
)

func TestChangelogQueryUnsupportedAttribue(t *testing.T) {
	_, err := newQueryEngine("publication_date", "json")
	require.Error(t, err)
}

func TestChangelogQueryUnsupportedTitle(t *testing.T) {
	_, err := newQueryEngine("title.size", "json")
	require.Error(t, err)
}

func TestChangelogQueryTitle(t *testing.T) {
	assertions := require.New(t)

	result, err := apply("title", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
	})
	assertions.NoError(err)
	assertions.Equal("changelog", result)
}

func TestChangelogQueryTitleAgainstRelease(t *testing.T) {
	assertions := require.New(t)

	result, err := apply("title", []changelog.Heading{
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
	})
	assertions.NoError(err)
	assertions.Empty(result)
}

func TestChangelogQueryUnsupportedEnter(t *testing.T) {
	assertions := require.New(t)

	query := &changelogQuery{}
	assertions.False(query.Enter(newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16")))
}

func TestChangelogQueryUnsupportedExit(t *testing.T) {
	assertions := require.New(t)

	query := &changelogQuery{}
	assertions.False(query.Exit(newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16")))
}

func TestChangelogQueryCollection(t *testing.T) {
	assertions := require.New(t)
	{
		query := &changelogQuery{}
		assertions.False(query.isCollection())
	}
	{
		query := &changelogQuery{projections{collection: true}}
		assertions.True(query.isCollection())
	}
}
