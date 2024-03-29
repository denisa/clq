package changelog

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAsPath(t *testing.T) {
	assert.Equal(t, "{}", asPath(""))
	assert.Equal(t, "{foo}", asPath("foo"))
}

func requireHeadingInterface(t *testing.T, expected string, actual Heading) {
	require.Equal(t, expected, actual.Title())
	require.Equal(t, asPath(actual.Title()), actual.String())
	require.Equal(t, actual.Title(), actual.DisplayTitle())
}
