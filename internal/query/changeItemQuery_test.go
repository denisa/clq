package query

import (
	"testing"

	"github.com/denisa/clq/internal/changelog"
	"github.com/stretchr/testify/require"
)

func TestChangeItemQueryMisformedSelector(t *testing.T) {
	_, err := newQueryEngine("releases[2].changes[].descriptions[", "json")
	require.Error(t, err)
}

func TestChangeItemQueryAsScalar(t *testing.T) {
	_, err := newQueryEngine("releases[2].changes[].descriptions", "json")
	require.Error(t, err)
}

func TestChangeItemQueryUnsupportedSelector(t *testing.T) {
	_, err := newQueryEngine("releases[2].changes[].descriptions[three]", "json")
	require.Error(t, err)
}

func TestChangeItemQueryUnsupportedAttribute(t *testing.T) {
	_, err := newQueryEngine("releases[2].changes[].descriptions[].fabulator", "json")
	require.Error(t, err)
}

func TestChangeItemQueryAgainstRelease(t *testing.T) {
	require := require.New(t)

	result, err := apply("releases[0].changes[].descriptions[]", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ChangeHeading, "Removed"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-16"),
	})
	require.NoError(err)
	require.JSONEq("[]", result)
}

func TestChangeQueryDescriptions(t *testing.T) {
	require := require.New(t)

	result, err := apply("releases[0].changes[].descriptions[]", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ChangeHeading, "Removed"),
		newHeading(changelog.ChangeDescription, "foo"),
	})
	require.NoError(err)
	require.Equal("[\"foo\"]", result)
}

func TestChangeItemQueryUnsupportedEnter(t *testing.T) {
	require := require.New(t)

	query := &changeItemQuery{}
	require.False(query.Enter(newHeading(changelog.IntroductionHeading, "changelog")))
}

func TestChangeItemQueryUnsupportedExit(t *testing.T) {
	require := require.New(t)

	query := &changeItemQuery{}
	require.False(query.Exit(newHeading(changelog.IntroductionHeading, "changelog")))
}

func TestChangeItemQueryCollection(t *testing.T) {
	require := require.New(t)
	{
		query := &changeItemQuery{}
		require.False(query.isCollection())
	}
	{
		query := &changeItemQuery{projections{collection: true}}
		require.True(query.isCollection())
	}
}
