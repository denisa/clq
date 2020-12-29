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
		changeMap ChangeMap
		expected  semver.Identifier
	}{
		{ChangeMap{"Added": true}, semver.Major},
		{ChangeMap{"Changed": true}, semver.Minor},
		{ChangeMap{"Fixed": true}, semver.Patch},
		{ChangeMap{"Unknown": true}, semver.Build},
		{ChangeMap{}, semver.Build},
	}
	for _, testcase := range testcases {
		t.Run(testcase.changeMap.String(), func(t *testing.T) {
			r := newChangeKind()
			require.Equal(t, testcase.expected, r.IncrementFor(testcase.changeMap))
		})
	}
}

func TestChangeKindFromUndefinedFile(t *testing.T) {
	_, err := NewChangeKind("testdata/undefined.json")
	require.Error(t, err)
}

func TestChangeKindFromWrongFileStructure(t *testing.T) {
	_, err := NewChangeKind("testdata/wrongFileStructure.json")
	require.Error(t, err)
}

func TestNewChangeKind(t *testing.T) {
	c, err := NewChangeKind("testdata/patch_only.json")
	require.NoError(t, err)
	require.Equal(t, "Fixed, Security", c.keysOf())
}
