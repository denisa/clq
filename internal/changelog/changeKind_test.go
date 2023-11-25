package changelog

import (
	"testing"

	"github.com/denisa/clq/internal/semver"
	"github.com/stretchr/testify/require"
)

func TestIsSupportedConfiguredValue(t *testing.T) {
	ck, _ := NewChangeKind("")
	_, err := ck.emojiFor("Fixed")
	require.NoError(t, err)
}

func TestIsSupportedUnconfiguredValue(t *testing.T) {
	ck, _ := NewChangeKind("")
	_, err := ck.emojiFor("Modified")
	require.Error(t, err)
}

func TestIsSupportedNoValue(t *testing.T) {
	ck, _ := NewChangeKind("")
	_, err := ck.emojiFor("")
	require.Error(t, err)
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
			ck, _ := NewChangeKind("")

			increment, trigger := ck.IncrementFor(testcase.changeMap)
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
	ck, err := NewChangeKind("testdata/patch_only.json")
	require.NoError(t, err)
	_, err = ck.emojiFor("unknown")
	require.Error(t, err)
}

func TestNewChangeKindWithoutEmoji(t *testing.T) {
	ck, err := NewChangeKind("testdata/patch_only.json")
	require.NoError(t, err)
	require.Equal(t, "Fixed, Security", ck.keysOf())
	{
		emoji, err := ck.emojiFor("Fixed")
		require.NoError(t, err)
		require.Equal(t, "", emoji)
	}
	{
		emoji, err := ck.emojiFor("Security")
		require.NoError(t, err)
		require.Equal(t, "", emoji)
	}
}

func TestNewChangeKindWithEmoji(t *testing.T) {
	ck, err := NewChangeKind("testdata/patch_only_with_emojis.json")
	require.NoError(t, err)
	require.Equal(t, "Fixed, Security", ck.keysOf())
	{
		emoji, err := ck.emojiFor("Fixed")
		require.NoError(t, err)
		require.Equal(t, "🐛", emoji)
	}
	{
		emoji, err := ck.emojiFor("Security")
		require.NoError(t, err)
		require.Equal(t, "🔒", emoji)
	}
}
