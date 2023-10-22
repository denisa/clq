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
	require := require.New(t)

	result, err := apply("releases[0].changes[]", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-16"),
	})
	require.NoError(err)
	require.JSONEq("[]", result)
}

func TestChangeQueryTitleSingle(t *testing.T) {
	require := require.New(t)

	result, err := apply("releases[0].changes[].title", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ChangeHeading, "Removed"),
	})
	require.NoError(err)
	require.JSONEq("[\"Removed\"]", result)
}

func TestChangeQueryTitleMultiple(t *testing.T) {
	require := require.New(t)

	result, err := apply("releases[0].changes[].title", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ChangeHeading, "Removed"),
		newHeading(changelog.ChangeHeading, "Added"),
		newHeading(changelog.ChangeHeading, "Security"),
	})
	require.NoError(err)
	require.JSONEq("[\"Removed\", \"Added\", \"Security\"]", result)
}

func TestChangeQueryWithoutItems(t *testing.T) {
	require := require.New(t)

	result, err := apply("releases[0].changes[]", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ChangeHeading, "Removed"),
	})
	require.NoError(err)
	require.JSONEq("[{\"title\":\"Removed\"}]", result)
}

func TestChangeQuerySecondReleaseChanges(t *testing.T) {
	require := require.New(t)

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
	require.NoError(err)
	require.JSONEq("[{\"title\":\"Added\"},{\"title\":\"Fixed\"},{\"title\":\"Security\"}]", result)
}

func TestChangeQuerySecondReleaseChangesRecursive(t *testing.T) {
	require := require.New(t)

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
	require.NoError(err)
	require.JSONEq("[{\"title\":\"Added\", \"descriptions\":[\"bar\"]},{\"title\":\"Fixed\", \"descriptions\":[\"waldo\"]}]", result)
}

func TestChangeQueryUnsupportedEnter(t *testing.T) {
	require := require.New(t)

	query := &changeQuery{}
	require.False(query.Enter(newHeading(changelog.IntroductionHeading, "changelog")))
}

func TestChangequeryUnsupportedExit(t *testing.T) {
	require := require.New(t)

	query := &changeQuery{}
	require.False(query.Exit(newHeading(changelog.IntroductionHeading, "changelog")))
}

func TestChangeQueryCollection(t *testing.T) {
	require := require.New(t)
	{
		query := &changeQuery{}
		require.False(query.isCollection())
	}
	{
		query := &changeQuery{projections{collection: true}}
		require.True(query.isCollection())
	}
}
