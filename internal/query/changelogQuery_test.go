package query

import (
	"testing"

	"github.com/denisa/clq/internal/changelog"
	"github.com/stretchr/testify/require"
)

func TestUnsupportedIntroductionQuery(t *testing.T) {
	_, err := NewQueryEngine("publication_date", "json")
	require.Error(t, err)
}

func TestUnsupportedIntroductionTitleQuery(t *testing.T) {
	_, err := NewQueryEngine("title.size", "json")
	require.Error(t, err)
}

func TestQueryTitle(t *testing.T) {
	require := require.New(t)

	result, err := apply("title", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
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

func TestUnsupportedChangelogEnter(t *testing.T) {
	require := require.New(t)

	query := &changelogQuery{}
	require.False(query.Enter(newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16")))
}

func TestUnsupportedChangelogExit(t *testing.T) {
	require := require.New(t)

	query := &changelogQuery{}
	require.False(query.Exit(newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16")))
}
