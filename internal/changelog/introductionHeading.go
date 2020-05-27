package changelog

import (
	"fmt"
)

// Level 1, introduction
type Introduction struct {
	title string
}

func newIntroduction(title string) (Heading, error) {
	if title == "" {
		return nil, fmt.Errorf("Validation error: Introduction’s title cannot stay empty")
	}
	return Introduction{title: title}, nil
}

func (h Introduction) Title() string {
	return h.title
}

func (h Introduction) String() string {
	return asPath(h.title)
}
