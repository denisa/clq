package changelog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnter(t *testing.T) {
	assertions := require.New(t)

	r := recorder{}
	assertions.Equal(len(r.events), 0)
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	h, _ := hf.newRelease("[1.2.4] - 2020-04-15")
	r.Enter(h)
	assertions.Equal(len(r.events), 1)
	assertions.Equal(r.events[0], "Enter "+h.String())
}

func TestExit(t *testing.T) {
	assertions := require.New(t)

	r := recorder{}
	assertions.Equal(len(r.events), 0)
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	h, _ := hf.newRelease("[1.2.4] - 2020-04-15")
	r.Exit(h)
	assertions.Equal(len(r.events), 1)
	assertions.Equal(r.events[0], "Exit "+h.String())
}
