package semver

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAllValues(t *testing.T) {
	assertions := require.New(t)
	_ = ForEach(func(i Identifier) error {
		t.Run(strconv.Itoa(int(i)), func(_ *testing.T) {
			val, err := NewIdentifier(i.String())
			assertions.NoError(err)
			assertions.Equal(i, val)
		})
		return nil
	})
}

func TestUndefinedValue(t *testing.T) {
	assertions := require.New(t)
	_, err := NewIdentifier("undefined value")
	assertions.Error(err)
	assertions.Panics(func() { _ = endOfEnum.String() })
}

func TestForEach(t *testing.T) {
	count := 0
	_ = ForEach(func(_ Identifier) error {
		count++
		return nil
	})
	require.Equal(t, int(endOfEnum-startOfEnum), count)
}

func TestForEachReturnsError(t *testing.T) {
	require.Error(t, ForEach(func(_ Identifier) error {
		return fmt.Errorf("dummy error")
	}))
}
