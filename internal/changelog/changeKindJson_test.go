package changelog

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshal(t *testing.T) {
	c := newChangeKind()
	val, err := json.Marshal(c)
	require.NoError(t, err)
	require.JSONEq(t, `[
    {"name":"Added", "increment":"major"},
    {"name":"Changed", "increment":"minor"},
    {"name":"Deprecated", "increment":"minor"},
    {"name":"Fixed", "increment":"patch"},
    {"name":"Removed", "increment":"major"},
    {"name":"Security", "increment":"patch"}
    ]`, string(val))
}

func TestUnmarshal(t *testing.T) {
	c := &ChangeKind{changes: make(changeKindToSemverIdentifier)}
	err := json.Unmarshal([]byte(`[
    {"name":"Added", "increment":"major"},
    {"name":"Fixed", "increment":"patch"}
    ]`), c)
	require.NoError(t, err)
	require.Equal(t, "Added, Fixed", c.keysOf())
}

func TestUnmarshalIllegalJson(t *testing.T) {
	c := &ChangeKind{changes: make(changeKindToSemverIdentifier)}
	err := json.Unmarshal([]byte(`{"name":"Added", "increment":"major"}`), c)
	require.Error(t, err)
}

func TestUnmarshalIllegalArgument(t *testing.T) {
	c := &ChangeKind{changes: make(changeKindToSemverIdentifier)}
	err := json.Unmarshal([]byte(`[
    {"name":"Added", "increment":"Major"},
    {"name":"Fixed", "increment":"Patch"}
    ]`), c)
	require.Error(t, err)
}
