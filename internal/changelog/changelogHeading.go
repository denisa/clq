package changelog

import (
	"fmt"
)

// Level 1, changelog
type Changelog struct {
	title string
}

func newChangelog(s string) (Heading, error) {
	if s == "" {
		return nil, fmt.Errorf("Validation error: Changelog cannot stay empty")
	}
	return Changelog{title: s}, nil
}

func (h Changelog) Name() string {
	return h.title
}

func (h Changelog) AsPath() string {
	return asPath(h.title)
}
