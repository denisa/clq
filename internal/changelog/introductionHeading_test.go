package changelog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewHeadingIntroduction(t *testing.T) {
	h, _ := NewHeading(TitleHeading, "changelog")
	requireHeadingInterface(t, "changelog", h)
}

func TestIntroduction(t *testing.T) {
	h, _ := newIntroduction("changelog")
	requireHeadingInterface(t, "changelog", h)
}

func TestEmptyIntroductionShouldFail(t *testing.T) {
	_, err := newIntroduction("")
	require.Error(t, err)
}
