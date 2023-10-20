package changelog

import (
	"testing"

	"github.com/denisa/clq/internal/semver"
	"github.com/stretchr/testify/require"
)

func TestIsSupportedConfiguredValue(t *testing.T) {
	r := newChangeKind()
	require.NoError(t, r.IsSupported("Fixed"))
}

func TestIsSupportedUnconfiguredValue(t *testing.T) {
	r := newChangeKind()
	require.Error(t, r.IsSupported("Modified"))
}

func TestIsSupportedNoValue(t *testing.T) {
	r := newChangeKind()
	require.Error(t, r.IsSupported(""))
}

func TestIncrementFor(t *testing.T) {
	testcases := []struct {
		changeMap         ChangeMap
		expectedIncrement semver.Identifier
		expectedTrigger   string
	}{
		{ChangeMap{"Added": true}, semver.Major, "Added"},
		{ChangeMap{"Changed": true}, semver.Minor, "Changed"},
		{ChangeMap{"Fixed": true}, semver.Patch, "Fixed"},
		{ChangeMap{"Unknown": true}, semver.Build, ""},
		{ChangeMap{}, semver.Build, ""},
	}
	for _, testcase := range testcases {
		t.Run(testcase.changeMap.String(), func(t *testing.T) {
			increment, trigger := newChangeKind().IncrementFor(testcase.changeMap)
			require.Equal(t, testcase.expectedIncrement, increment)
			require.Equal(t, testcase.expectedTrigger, trigger)
		})
	}
}

func TestChangeKindFromUndefinedFile(t *testing.T) {
	_, err := NewChangeKind("testdata/undefined.json")
	require.Error(t, err)
}

func TestChangeKindNameMissing(t *testing.T) {
	_, err := NewChangeKind("testdata/missing_name.json")
	require.Error(t, err)
}

func TestChangeKindIncrementMissing(t *testing.T) {
	_, err := NewChangeKind("testdata/missing_increment.json")
	require.Error(t, err)
}

func TestChangeKindFromWrongFileStructure(t *testing.T) {
	_, err := NewChangeKind("testdata/wrongFileStructure.json")
	require.Error(t, err)
}

func TestEmojiWrongHeading(t *testing.T) {
	c, err := NewChangeKind("testdata/patch_only.json")
	require.NoError(t, err)
	_, err = c.emojiFor("unknown")
	require.Error(t, err)
}

func TestNewChangeKindWithoutEmoji(t *testing.T) {
	c, err := NewChangeKind("testdata/patch_only.json")
	require.NoError(t, err)
	require.Equal(t, "Fixed, Security", c.keysOf())
	{
		emoji, err := c.emojiFor("Fixed")
		require.NoError(t, err)
		require.Equal(t, "", emoji)
	}
	{
		emoji, err := c.emojiFor("Security")
		require.NoError(t, err)
		require.Equal(t, "", emoji)
	}
}

func TestNewChangeKindWithEmoji(t *testing.T) {
	c, err := NewChangeKind("testdata/patch_only_with_emojis.json")
	require.NoError(t, err)
	require.Equal(t, "Fixed, Security", c.keysOf())
	{
		emoji, err := c.emojiFor("Fixed")
		require.NoError(t, err)
		require.Equal(t, "🐛", emoji)
	}
	{
		emoji, err := c.emojiFor("Security")
		require.NoError(t, err)
		require.Equal(t, "🔒", emoji)
	}
}
