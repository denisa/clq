package changelog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnter(t *testing.T) {
	require := require.New(t)

	r := recorder{}
	require.Equal(len(r.events), 0)
	h, _ := newRelease("[1.2.4] - 2020-04-15")
	r.Enter(h)
	require.Equal(len(r.events), 1)
	require.Equal(r.events[0], "Enter "+h.String())
}

func TestExit(t *testing.T) {
	require := require.New(t)

	r := recorder{}
	require.Equal(len(r.events), 0)
	h, _ := newRelease("[1.2.4] - 2020-04-15")
	r.Exit(h)
	require.Equal(len(r.events), 1)
	require.Equal(r.events[0], "Exit "+h.String())
}
