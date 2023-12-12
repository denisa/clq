package changelog

import "strings"

// ChangeMap tracks the Changes that have been defined in a release
type ChangeMap map[string]bool

func (c ChangeMap) String() string {
	var changes []string
	for k := range c {
		changes = append(changes, k)
	}
	return strings.Join(changes, ", ")
}
