package changelog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewHeadingChange(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	h, _ := hf.NewHeading(ChangeHeading, "Security")
	requireHeadingInterface(t, "Security", h)
}

func TestChange(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	h, _ := hf.newChange("Security")
	requireHeadingInterface(t, "Security", h)
}

func TestEmptyChangeShouldFail(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	_, err := hf.newChange("")
	require.Error(t, err)
}

func TestChangeDisplayTitleWithEmoji(t *testing.T) {
	ck, _ := NewChangeKind("testdata/patch_only_with_emojis.json")
	hf := NewHeadingFactory(ck)
	h, _ := hf.newChange("Security")
	require.Equal(t, "Security", h.Title())
	require.Equal(t, "🔒 Security", h.DisplayTitle())
}
