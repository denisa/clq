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

func TestUnknownChangeShouldFail(t *testing.T) {
	_, err := newChange("No-op")
	require.Error(t, err)
}
