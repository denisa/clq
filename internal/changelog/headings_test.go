package changelog

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAsPath(t *testing.T) {
	assert.Equal(t, "{}", asPath(""))
	assert.Equal(t, "{foo}", asPath("foo"))
}

func TestNewHeadingUnknown(t *testing.T) {
	_, err := NewHeading(-1, "Who knows what")
	require.Error(t, err)
}

func requireHeadingInterface(t *testing.T, name string, actual Heading) {
	require.Equal(t, name, actual.Name())
	require.Equal(t, asPath(actual.Name()), actual.AsPath())
}

func requireHeadingEquals(require *require.Assertions, expected Heading, actual Heading) {
	require.Equal(expected, actual)
	require.Equal(reflect.TypeOf(expected), reflect.TypeOf(actual))
}
