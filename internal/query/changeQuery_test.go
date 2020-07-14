package query

import (
	"testing"

	"github.com/denisa/clq/internal/changelog"
	"github.com/stretchr/testify/require"
)

func TestUnsupportedChangeQuerySelector(t *testing.T) {
	_, err := NewQueryEngine("release[2].changes[3]")
	require.Error(t, err)
}

func TestUnsupportedChangeQueryAttribute(t *testing.T) {
	_, err := NewQueryEngine("release[2].changes[].title")
	require.Error(t, err)
}

func TestQueryChangesWithoutItems(t *testing.T) {
	require := require.New(t)

	result, err := apply("releases[0].changes[]", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ChangeHeading, "Removed"),
	})
	require.NoError(err)
	require.Equal("{\"kind\":\"Removed\", \"changes\":[]}", result)
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
	})
	require.NoError(err)
	require.Equal("{\"kind\":\"Added\", \"changes\":[\"bar\",]}", result)
}
