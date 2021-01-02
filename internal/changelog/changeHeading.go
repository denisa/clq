package changelog

import "fmt"

// Level 3, change groups
type Change struct {
	heading
}

func newChange(title string) (Heading, error) {
	if title == "" {
		return nil, fmt.Errorf("Validation error: change cannot stay empty")
	}
	return Change{heading{title: title, kind: ChangeHeading}}, nil
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
