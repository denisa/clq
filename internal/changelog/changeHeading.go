package changelog

import (
	"fmt"
	"regexp"
)

// Level 3, change groups
type Change struct {
	title string
}

func newChange(title string) (Heading, error) {
	for val := range changeKind {
		if matched, _ := regexp.MatchString(`^`+val+`$`, title); matched {
			return Change{title: title}, nil
		}
	}

	return nil, fmt.Errorf("Validation error: Unknown change headings %q is not one of [%v]", title, keysOf(changeKind))
}

func (h Change) Title() string {
	return h.title
}

func (h Change) String() string {
	return asPath(h.title)
}
