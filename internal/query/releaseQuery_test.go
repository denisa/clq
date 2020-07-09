package query

import (
	"testing"

	"github.com/denisa/clq/internal/changelog"
	"github.com/stretchr/testify/require"
)

func TestMisformedReleaseQuerySelector(t *testing.T) {
	_, err := NewQueryEngine("releases[", "json")
	require.Error(t, err)
}

func TestUnsupportedReleaseQuerySelector(t *testing.T) {
	_, err := NewQueryEngine("releases[three]", "json")
	require.Error(t, err)
}

func TestUnsupportedReleaseQuery(t *testing.T) {
	_, err := NewQueryEngine("releases[3].publication_date", "json")
	require.Error(t, err)
}

func TestUnsupportedReleaseTitleQuery(t *testing.T) {
	_, err := NewQueryEngine("releases[2].title.size", "json")
	require.Error(t, err)
}

func TestUnsupportedReleaseDateQuery(t *testing.T) {
	_, err := NewQueryEngine("releases[2].date.size", "json")
	require.Error(t, err)
}

func TestUnsupportedReleaseLabelQuery(t *testing.T) {
	_, err := NewQueryEngine("releases[2].label.size", "json")
	require.Error(t, err)
}

func TestUnsupportedReleaseStatusQuery(t *testing.T) {
	_, err := NewQueryEngine("releases[2].status.size", "json")
	require.Error(t, err)
}

func TestUnsupportedReleaseVersionQuery(t *testing.T) {
	_, err := NewQueryEngine("releases[2].version.size", "json")
	require.Error(t, err)
}

func TestUnsupportedReleaseIndexQuery(t *testing.T) {
	_, err := NewQueryEngine("releases[three].date", "json")
	require.Error(t, err)
}

func TestQueryReleaseAgainstIntroduction(t *testing.T) {
	require := require.New(t)

	result, err := apply("releases[1].version", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.IntroductionHeading, "there should not be another level 1 title"),
	})
	require.NoError(err)
	require.Empty(result)
}

func TestQuerySecondReleaseVersion(t *testing.T) {
	require := require.New(t)

	result, err := apply("releases[1].version", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-15 Cabrel"),
	})
	require.NoError(err)
	require.Equal("1.2.2", result)
}

func TestQuerySecondReleaseDate(t *testing.T) {
	require := require.New(t)

	result, err := apply("releases[0].date", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-15 Cabrel"),
	})
	require.NoError(err)
	require.Equal("2020-05-16", result)
}

func TestQuerySecondReleaseLabel(t *testing.T) {
	require := require.New(t)

	result, err := apply("releases[1].label", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-15 Cabrel"),
	})
	require.NoError(err)
	require.Equal("Cabrel", result)
}

func TestQuerySecondReleaseStatusPrereleased(t *testing.T) {
	require := require.New(t)

	result, err := apply("releases[1].status", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ReleaseHeading, "[1.2.2-rc.1] - 2020-05-15 Cabrel"),
	})
	require.NoError(err)
	require.Equal("prereleased", result)
}

func TestQuerySecondReleaseStatusReleased(t *testing.T) {
	require := require.New(t)

	result, err := apply("releases[1].status", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-15 Cabrel"),
	})
	require.NoError(err)
	require.Equal("released", result)
}

func TestQueryFirstReleaseStatusUnreleased(t *testing.T) {
	require := require.New(t)

	result, err := apply("releases[0].status", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[Unreleased]"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-15 Cabrel"),
	})
	require.NoError(err)
	require.Equal("unreleased", result)
}

func TestQuerySecondReleaseStatusYanked(t *testing.T) {
	require := require.New(t)

	result, err := apply("releases[1].status", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ReleaseHeading, "1.2.2 - 2020-05-15 [YANKED]"),
	})
	require.NoError(err)
	require.Equal("yanked", result)
}

func TestQuerySecondRelease(t *testing.T) {
	require := require.New(t)

	result, err := apply("releases[1]", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-15 Cabrel"),
	})
	require.NoError(err)
	require.JSONEq("{\"version\":\"1.2.2\", \"date\":\"2020-05-15\"}", result)
}

func TestUnsupportedReleaseEnter(t *testing.T) {
	require := require.New(t)

	query := &releaseQuery{}
	require.False(query.Enter(newHeading(changelog.IntroductionHeading, "changelog")))
}

func TestUnsupportedReleaseItemExit(t *testing.T) {
	require := require.New(t)

	query := &releaseQuery{}
	require.False(query.Exit(newHeading(changelog.IntroductionHeading, "changelog")))
}
