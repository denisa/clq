package changelog

import (
	"fmt"
)

// Level 1, introduction
type Introduction struct {
	title string
}

func newIntroduction(s string) (Heading, error) {
	if s == "" {
		return nil, fmt.Errorf("Validation error: Introduction’s title cannot stay empty")
	}
	return Introduction{title: s}, nil
}

func (h Introduction) Name() string {
	return h.title
}

func (h Introduction) String() string {
	return asPath(h.title)
}
