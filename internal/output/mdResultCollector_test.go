package output

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMdNoOutputDefined(t *testing.T) {
	require.Equal(t, "", formatNoOutputDefined("md"))
}

func TestMdIntroductionHeading(t *testing.T) {
	require.Equal(t, "# Changelog", formatIntroductionHeading("md"))
}

func TestMdReleaseHeading(t *testing.T) {
	require.Equal(t, "## [1.2.3] - 2020-05-16", formatReleaseHeading("md"))
}

func TestMdChangeHeading(t *testing.T) {
	require.Equal(t, "### Added", formatChangeHeading("md"))
}

func TestMdChangeDescription(t *testing.T) {
	require.Equal(t, "- foo", formatChangeDescription("md"))
}

func TestMdLoneArray(t *testing.T) {
	require.Equal(t, "- foo\n- bar", formatLoneArray("md"))
}
