package changelog

import (
	"testing"
)

func TestNewHeadingChangeDescription(t *testing.T) {
	h, _ := NewHeading(ChangeDescription, "foo")
	requireHeadingInterface(t, "foo", h)
}

func TestChangeDescription(t *testing.T) {
	h, _ := newChangeItem("foo")
	requireHeadingInterface(t, "foo", h)
}
