package changelog

import (
	"fmt"
)

// HeadingKind is the type for the multiple sections.
type HeadingKind int

const (
	IntroductionHeading HeadingKind = iota
	ReleaseHeading
	ChangeHeading
	ChangeDescription
)

// A Heading is the interface common to every sections.
type Heading interface {
	// The Title of the section
	Title() string
	// The HeadingKind of the section
	Kind() HeadingKind
	String() string
}

type heading struct {
	title string
	kind  HeadingKind
}

func asPath(name string) string {
	return "{" + name + "}"
}

// NewHeading is the factory method that, given a kind and a title, returns the appropriate Heading.
func NewHeading(kind HeadingKind, title string) (Heading, error) {
	switch kind {
	case IntroductionHeading:
		return newIntroduction(title)
	case ReleaseHeading:
		return newRelease(title)
	case ChangeHeading:
		return newChange(title)
	case ChangeDescription:
		return newChangeItem(title)
	}
	return nil, fmt.Errorf("Unknown heading kind %v", kind)
}
