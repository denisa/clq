package query

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJsonNoOutputDefined(t *testing.T) {
	require.Equal(t, "{}", formatNoOutputDefined("json"))
}

func TestJsonIntroductionHeading(t *testing.T) {
	require.Equal(t, "{\"title\":\"Changelog\"}", formatIntroductionHeading("json"))
}

func TestJsonReleaseHeading(t *testing.T) {
	require.Equal(t, "{\"title\":\"[1.2.3] - 2020-05-16\"}", formatReleaseHeading("json"))
}

func TestJsonChangeHeading(t *testing.T) {
	require.Equal(t, "{\"title\":\"Added\"}", formatChangeHeading("json"))
}

func TestJsonChangeDescription(t *testing.T) {
	require.Equal(t, "foo", formatChangeDescription("json"))
}

func TestJsonLoneArray(t *testing.T) {
	require.Equal(t, "{\"changes\":[\"foo\",\"bar\"]}", formatLoneArray("json"))
}
