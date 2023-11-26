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
	assertions := require.New(t)

	result, err := apply("releases[0].changes[].descriptions[]", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ChangeHeading, "Removed"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-16"),
	})
	assertions.NoError(err)
	assertions.JSONEq("[]", result)
}

func TestChangeQueryDescriptions(t *testing.T) {
	assertions := require.New(t)

	result, err := apply("releases[0].changes[].descriptions[]", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ChangeHeading, "Removed"),
		newHeading(changelog.ChangeDescription, "foo"),
	})
	assertions.NoError(err)
	assertions.Equal("[\"foo\"]", result)
}

func TestChangeItemQueryUnsupportedEnter(t *testing.T) {
	assertions := require.New(t)

	query := &changeItemQuery{}
	assertions.False(query.Enter(newHeading(changelog.IntroductionHeading, "changelog")))
}

func TestChangeItemQueryUnsupportedExit(t *testing.T) {
	assertions := require.New(t)

	query := &changeItemQuery{}
	assertions.False(query.Exit(newHeading(changelog.IntroductionHeading, "changelog")))
}

func TestChangeItemQueryCollection(t *testing.T) {
	assertions := require.New(t)
	{
		query := &changeItemQuery{}
		assertions.False(query.isCollection())
	}
	{
		query := &changeItemQuery{projections{collection: true}}
		assertions.True(query.isCollection())
	}
}
