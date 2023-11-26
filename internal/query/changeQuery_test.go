package query

import (
	"testing"

	"github.com/denisa/clq/internal/changelog"
	"github.com/stretchr/testify/require"
)

func TestChangeQueryMisformedSelector(t *testing.T) {
	_, err := newQueryEngine("releases[2].changes[", "json")
	require.Error(t, err)
}

func TestChangeQueryAsScalar(t *testing.T) {
	_, err := newQueryEngine("releases[2].changes", "json")
	require.Error(t, err)
}

func TestChangeQueryUnsupportedSelector(t *testing.T) {
	_, err := newQueryEngine("releases[2].changes[three]", "json")
	require.Error(t, err)
}

func TestChangeQueryUnsupportedAttribute(t *testing.T) {
	_, err := newQueryEngine("releases[2].changes[].fabulator", "json")
	require.Error(t, err)
}

func TestChangeQueryUnsupportedTitle(t *testing.T) {
	_, err := newQueryEngine("releases[2].changes[].title.size", "json")
	require.Error(t, err)
}

func TestChangeQueryAgainstIntroduction(t *testing.T) {
	assertions := require.New(t)

	result, err := apply("releases[0].changes[]", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-16"),
	})
	assertions.NoError(err)
	assertions.JSONEq("[]", result)
}

func TestChangeQueryTitleSingle(t *testing.T) {
	assertions := require.New(t)

	result, err := apply("releases[0].changes[].title", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ChangeHeading, "Removed"),
	})
	assertions.NoError(err)
	assertions.JSONEq("[\"Removed\"]", result)
}

func TestChangeQueryTitleMultiple(t *testing.T) {
	assertions := require.New(t)

	result, err := apply("releases[0].changes[].title", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ChangeHeading, "Removed"),
		newHeading(changelog.ChangeHeading, "Added"),
		newHeading(changelog.ChangeHeading, "Security"),
	})
	assertions.NoError(err)
	assertions.JSONEq("[\"Removed\", \"Added\", \"Security\"]", result)
}

func TestChangeQueryWithoutItems(t *testing.T) {
	assertions := require.New(t)

	result, err := apply("releases[0].changes[]", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ChangeHeading, "Removed"),
	})
	assertions.NoError(err)
	assertions.JSONEq("[{\"title\":\"Removed\"}]", result)
}

func TestChangeQuerySecondReleaseChanges(t *testing.T) {
	assertions := require.New(t)

	result, err := apply("releases[1].changes[]", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ChangeHeading, "Removed"),
		newHeading(changelog.ChangeDescription, "foo"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-15 Cabrel"),
		newHeading(changelog.ChangeHeading, "Added"),
		newHeading(changelog.ChangeDescription, "bar"),
		newHeading(changelog.ChangeHeading, "Fixed"),
		newHeading(changelog.ChangeDescription, "waldo"),
		newHeading(changelog.ChangeHeading, "Security"),
		newHeading(changelog.ChangeDescription, "thud"),
	})
	assertions.NoError(err)
	assertions.JSONEq("[{\"title\":\"Added\"},{\"title\":\"Fixed\"},{\"title\":\"Security\"}]", result)
}

func TestChangeQuerySecondReleaseChangesRecursive(t *testing.T) {
	assertions := require.New(t)

	result, err := apply("releases[1].changes[]/", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ChangeHeading, "Removed"),
		newHeading(changelog.ChangeDescription, "foo"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-15 Cabrel"),
		newHeading(changelog.ChangeHeading, "Added"),
		newHeading(changelog.ChangeDescription, "bar"),
		newHeading(changelog.ChangeHeading, "Fixed"),
		newHeading(changelog.ChangeDescription, "waldo"),
	})
	assertions.NoError(err)
	assertions.JSONEq("[{\"title\":\"Added\", \"descriptions\":[\"bar\"]},{\"title\":\"Fixed\", \"descriptions\":[\"waldo\"]}]", result)
}

func TestChangeQueryUnsupportedEnter(t *testing.T) {
	assertions := require.New(t)

	query := &changeQuery{}
	assertions.False(query.Enter(newHeading(changelog.IntroductionHeading, "changelog")))
}

func TestChangequeryUnsupportedExit(t *testing.T) {
	assertions := require.New(t)

	query := &changeQuery{}
	assertions.False(query.Exit(newHeading(changelog.IntroductionHeading, "changelog")))
}

func TestChangeQueryCollection(t *testing.T) {
	assertions := require.New(t)
	{
		query := &changeQuery{}
		assertions.False(query.isCollection())
	}
	{
		query := &changeQuery{projections{collection: true}}
		assertions.True(query.isCollection())
	}
}
