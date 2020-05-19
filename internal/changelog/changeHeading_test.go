package changelog

import (
	"testing"
)

func TestNewHeadingChange(t *testing.T) {
	h, _ := NewHeading(ChangeHeading, "Security")
	requireHeadingInterface(t, "Security", h)
}

func TestChange(t *testing.T) {
	h, _ := newChange("Security")
	requireHeadingInterface(t, "Security", h)
}
