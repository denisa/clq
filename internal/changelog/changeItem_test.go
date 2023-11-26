package changelog

import (
	"testing"
)

func TestNewHeadingChangeDescription(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	h, _ := hf.NewHeading(ChangeDescription, "foo")
	requireHeadingInterface(t, "foo", h)
}

func TestChangeDescription(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	h, _ := hf.newChangeItem("foo")
	requireHeadingInterface(t, "foo", h)
}
