package changelog

import (
	"fmt"
)

// Level 1, introduction
type Introduction struct {
	heading
}

func (h HeadingsFactory) newIntroduction(title string) (Heading, error) {
	if title == "" {
		return nil, fmt.Errorf("Validation error: Introduction’s title cannot stay empty")
	}
	return Introduction{heading{title: title, kind: IntroductionHeading}}, nil
}

func (h Introduction) DisplayTitle() string {
	return h.Title()
}

func (h Introduction) Title() string {
	return h.title
}

func (h Introduction) Kind() HeadingKind {
	return h.kind
}

func (h Introduction) String() string {
	return asPath(h.title)
}
