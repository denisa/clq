package changelog

import (
	"fmt"
	"regexp"
)

// Level 3, change groups
type Change struct {
	heading
}

func newChange(title string) (Heading, error) {
	for val := range changeKind {
		if matched, _ := regexp.MatchString(`^`+val+`$`, title); matched {
			return Change{heading{title: title, kind: ChangeHeading}}, nil
		}
	}

	return nil, fmt.Errorf("Validation error: Unknown change headings %q is not one of [%v]", title, keysOf(changeKind))
}

func (h Change) Title() string {
	return h.title
}

func (h Change) Kind() HeadingKind {
	return h.kind
}

func (h Change) String() string {
	return asPath(h.title)
}
