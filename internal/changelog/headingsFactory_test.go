package changelog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewHeadingUnknown(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	_, err := hf.NewHeading(-1, "Who knows what")
	require.Error(t, err)
}
