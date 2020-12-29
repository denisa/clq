package changelog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewHeadingChange(t *testing.T) {
	h, _ := NewHeading(ChangeHeading, "Security")
	requireHeadingInterface(t, "Security", h)
}

func TestChange(t *testing.T) {
	h, _ := newChange("Security")
	requireHeadingInterface(t, "Security", h)
}

func TestEmptyChangeShouldFail(t *testing.T) {
	_, err := newChange("")
	require.Error(t, err)
}
