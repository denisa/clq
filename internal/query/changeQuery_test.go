package query

import (
	"testing"

	"github.com/denisa/clq/internal/changelog"
	"github.com/stretchr/testify/require"
)

func TestMisformedChangeQuerySelector(t *testing.T) {
	_, err := NewQueryEngine("releases[2].changes[", "json")
	require.Error(t, err)
}

func TestUnsupportedChangeQuerySelector(t *testing.T) {
	_, err := NewQueryEngine("releases[2].changes[three]", "json")
	require.Error(t, err)
}

func TestUnsupportedChangeQueryAttribute(t *testing.T) {
	_, err := NewQueryEngine("releases[2].changes[].fabulator", "json")
	require.Error(t, err)
}

func TestUnsupportedChangeTitleQuery(t *testing.T) {
	_, err := NewQueryEngine("releases[2].changes[].title.size", "json")
	require.Error(t, err)
}

func TestQueryChangeAgainstIntroduction(t *testing.T) {
	require := require.New(t)

	result, err := apply("releases[0].changes[]", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-16"),
	})
	require.NoError(err)
	require.Empty(result)
}

func TestQueryChangeTitle(t *testing.T) {
	require := require.New(t)

	result, err := apply("releases[0].changes[].title", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ChangeHeading, "Removed"),
	})
	require.NoError(err)
	require.Equal("Removed", result)
}

func TestQueryChangesWithoutItems(t *testing.T) {
	require := require.New(t)

	result, err := apply("releases[0].changes[]", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ChangeHeading, "Removed"),
	})
	require.NoError(err)
	require.JSONEq("{\"title\":\"Removed\"}", result)
}

func TestQuerySecondReleaseChanges(t *testing.T) {
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

func TestQuerySecondReleaseChangesRecursive(t *testing.T) {
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
	require.JSONEq("[{\"title\":\"Added\", \"changes\":[\"bar\"]},{\"title\":\"Fixed\", \"changes\":[\"waldo\"]}]", result)
}

func TestUnsupportedChangeEnter(t *testing.T) {
	require := require.New(t)

	query := &changeQuery{}
	require.False(query.Enter(newHeading(changelog.IntroductionHeading, "changelog")))
}

func TestUnsupportedChangeExit(t *testing.T) {
	require := require.New(t)

	query := &changeQuery{}
	require.False(query.Exit(newHeading(changelog.IntroductionHeading, "changelog")))
}
