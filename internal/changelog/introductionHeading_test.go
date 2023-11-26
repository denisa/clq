package changelog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewHeadingIntroduction(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	h, _ := hf.NewHeading(IntroductionHeading, "changelog")
	requireHeadingInterface(t, "changelog", h)
}

func TestIntroduction(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	h, _ := hf.newIntroduction("changelog")
	requireHeadingInterface(t, "changelog", h)
}

func TestEmptyIntroductionShouldFail(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	_, err := hf.newIntroduction("")
	require.Error(t, err)
}
