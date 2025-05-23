package query

import (
	"testing"

	"github.com/denisa/clq/internal/changelog"
	"github.com/stretchr/testify/require"
)

func TestReleaseQueryMisformedSelector(t *testing.T) {
	_, err := newQueryEngine("releases[", "json")
	require.Error(t, err)
}

func TestReleaseQueryUnsupportedSelector(t *testing.T) {
	_, err := newQueryEngine("releases[three]", "json")
	require.Error(t, err)
}

func TestReleaseQueryAsScalar(t *testing.T) {
	_, err := newQueryEngine("releases", "json")
	require.Error(t, err)
}

func TestReleaseQueryUnsupportedAttribute(t *testing.T) {
	_, err := newQueryEngine("releases[3].publication_date", "json")
	require.Error(t, err)
}

func TestReleaseQueryUnsupportedTitle(t *testing.T) {
	_, err := newQueryEngine("releases[2].title.size", "json")
	require.Error(t, err)
}

func TestReleaseQueryUnsupportedDate(t *testing.T) {
	_, err := newQueryEngine("releases[2].date.size", "json")
	require.Error(t, err)
}

func TestReleaseQueryUnsupportedLabel(t *testing.T) {
	_, err := newQueryEngine("releases[2].label.size", "json")
	require.Error(t, err)
}

func TestReleaseQueryUnsupportedStatus(t *testing.T) {
	_, err := newQueryEngine("releases[2].status.size", "json")
	require.Error(t, err)
}

func TestReleaseQueryUnsupportedVersion(t *testing.T) {
	_, err := newQueryEngine("releases[2].version.size", "json")
	require.Error(t, err)
}

func TestReleaseQueryUnsupportedIndex(t *testing.T) {
	_, err := newQueryEngine("releases[three].date", "json")
	require.Error(t, err)
}

func TestReleaseQueryAgainstIntroduction(t *testing.T) {
	assertions := require.New(t)

	result, err := apply("releases[1].version", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.IntroductionHeading, "there should not be another level 1 title"),
	})
	assertions.NoError(err)
	assertions.Empty(result)
}

func TestReleaseQuerySecondReleaseVersion(t *testing.T) {
	assertions := require.New(t)

	result, err := apply("releases[1].version", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-15 Cabrel"),
	})
	assertions.NoError(err)
	assertions.Equal("1.2.2", result)
}

func TestReleaseQuerySecondReleaseDate(t *testing.T) {
	assertions := require.New(t)

	result, err := apply("releases[0].date", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-15 Cabrel"),
	})
	assertions.NoError(err)
	assertions.Equal("2020-05-16", result)
}

func TestReleaseQuerySecondReleaseLabel(t *testing.T) {
	assertions := require.New(t)

	result, err := apply("releases[1].label", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-15 Cabrel"),
	})
	assertions.NoError(err)
	assertions.Equal("Cabrel", result)
}

func TestReleaseQuerySecondReleaseStatusPrereleased(t *testing.T) {
	assertions := require.New(t)

	result, err := apply("releases[1].status", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ReleaseHeading, "[1.2.2-rc.1] - 2020-05-15 Cabrel"),
	})
	assertions.NoError(err)
	assertions.Equal("prereleased", result)
}

func TestReleaseQuerySecondReleaseStatusReleased(t *testing.T) {
	assertions := require.New(t)

	result, err := apply("releases[1].status", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-15 Cabrel"),
	})
	assertions.NoError(err)
	assertions.Equal("released", result)
}

func TestReleaseQueryFirstReleaseStatusUnreleased(t *testing.T) {
	assertions := require.New(t)

	result, err := apply("releases[0].status", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[Unreleased]"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-15 Cabrel"),
	})
	assertions.NoError(err)
	assertions.Equal("unreleased", result)
}

func TestReleaseQueryFirstReleaseDateUnreleased(t *testing.T) {
	assertions := require.New(t)

	result, err := apply("releases[0].date", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[Unreleased]"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-15 Cabrel"),
	})
	assertions.NoError(err)
	assertions.Equal("", result)
}

func TestReleaseQuerySecondReleaseStatusYanked(t *testing.T) {
	assertions := require.New(t)

	result, err := apply("releases[1].status", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-15 [YANKED]"),
	})
	assertions.NoError(err)
	assertions.Equal("yanked", result)
}

func TestReleaseQuerySecondReleaseTitle(t *testing.T) {
	assertions := require.New(t)

	result, err := apply("releases[1].title", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-15 [YANKED]"),
	})
	assertions.NoError(err)
	assertions.Equal("[1.2.2] - 2020-05-15 [YANKED]", result)
}

func TestReleaseQuerySecondRelease(t *testing.T) {
	assertions := require.New(t)

	result, err := apply("releases[1]", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-15 Cabrel"),
	})
	assertions.NoError(err)
	assertions.JSONEq("{\"version\":\"1.2.2\", \"date\":\"2020-05-15\"}", result)
}

func TestReleaseQueryUnsupportedEnter(t *testing.T) {
	assertions := require.New(t)

	query := &releaseQuery{}
	assertions.False(query.Enter(newHeading(changelog.IntroductionHeading, "changelog")))
}

func TestReleaseQueryUnsupportedRExit(t *testing.T) {
	assertions := require.New(t)

	query := &releaseQuery{}
	assertions.False(query.Exit(newHeading(changelog.IntroductionHeading, "changelog")))
}

func TestReleaseQueryCollection(t *testing.T) {
	assertions := require.New(t)
	{
		query := &releaseQuery{}
		assertions.False(query.isCollection())
	}
	{
		query := &releaseQuery{}
		query.collection = true
		assertions.True(query.isCollection())
	}
}
