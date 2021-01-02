package semver

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAllValues(t *testing.T) {
	require := require.New(t)
	ForEach(func(i Identifier) error {
		t.Run(strconv.Itoa(int(i)), func(t *testing.T) {
			val, err := NewIdentifier(i.String())
			require.NoError(err)
			require.Equal(i, val)
		})
		return nil
	})
}

func TestUndefinedValue(t *testing.T) {
	require := require.New(t)
	_, err := NewIdentifier("undefined value")
	require.Error(err)
	require.Panics(func() { endOfEnum.String() })
}

func TestForEach(t *testing.T) {
	count := 0
	ForEach(func(arg1 Identifier) error {
		count++
		return nil
	})
	require.Equal(t, int(endOfEnum-startOfEnum), count)
}

func TestForEachReturnsError(t *testing.T) {
	require.Error(t, ForEach(func(arg1 Identifier) error {
		return fmt.Errorf("Dummy error")
	}))
}
