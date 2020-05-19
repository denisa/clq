package changelog

import (
	"fmt"
	"regexp"
)

// Level 3, change groups
type Change struct {
	name string
}

func newChange(s string) (Heading, error) {
	for val := range changeKind {
		if matched, _ := regexp.MatchString(`^`+val+`$`, s); matched {
			return Change{name: s}, nil
		}
	}

	return nil, fmt.Errorf("Validation error: Unknown change headings %q is not one of [%v]", s, keysOf(changeKind))
}

func (h Change) Name() string {
	return h.name
}
func (h Change) AsPath() string {
	return asPath(h.name)
}
