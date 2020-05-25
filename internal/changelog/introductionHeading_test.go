package changelog

import (
	"testing"
)

func TestNewHeadingIntroduction(t *testing.T) {
	h, _ := NewHeading(TitleHeading, "changelog")
	requireHeadingInterface(t, "changelog", h)
}

func TestIntroduction(t *testing.T) {
	h, _ := newIntroduction("changelog")
	requireHeadingInterface(t, "changelog", h)
}
