package changelog

import (
	"testing"
)

func TestNewHeadingChangelog(t *testing.T) {
	h, _ := NewHeading(TitleHeading, "changelog")
	requireHeadingInterface(t, "changelog", h)
}

func TestChangelog(t *testing.T) {
	h, _ := newChangelog("changelog")
	requireHeadingInterface(t, "changelog", h)
}
