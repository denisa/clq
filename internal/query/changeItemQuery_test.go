package query

import (
	"testing"

	"github.com/denisa/clq/internal/changelog"
	"github.com/stretchr/testify/require"
)

func TestUnsupportedChangeItemQuerySelector(t *testing.T) {
	_, err := NewQueryEngine("releases[2].changes[].descriptions[three]")
	require.Error(t, err)
}

func TestUnsupportedChangeItemQueryAttribute(t *testing.T) {
	_, err := NewQueryEngine("releases[2].changes[].descriptions[].fabulator")
	require.Error(t, err)
}

func TestQueryChangeItemAgainstRelease(t *testing.T) {
	require := require.New(t)

	result, err := apply("releases[0].changes[].descriptions[]", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ChangeHeading, "Removed"),
		newHeading(changelog.ReleaseHeading, "[1.2.2] - 2020-05-16"),
	})
	require.NoError(err)
	require.Empty(result)
}

func TestQueryChangeDescriptions(t *testing.T) {
	require := require.New(t)

	result, err := apply("releases[0].changes[].descriptions[]", []changelog.Heading{
		newHeading(changelog.IntroductionHeading, "changelog"),
		newHeading(changelog.ReleaseHeading, "[1.2.3] - 2020-05-16"),
		newHeading(changelog.ChangeHeading, "Removed"),
		newHeading(changelog.ChangeDescription, "foo"),
	})
	require.NoError(err)
	require.Equal("foo", result)
}
