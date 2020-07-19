package changelog

import (
	"fmt"
	"sort"
	"strings"
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

type incrementKind int

const (
	semverMajor incrementKind = iota
	semverMinor
	semverPatch
	semverPrerelease
	semverBuild
)

type changeToIncrementKind map[string]incrementKind

var changeKind = changeToIncrementKind{"Added": semverMajor, "Removed": semverMajor, "Changed": semverMinor, "Deprecated": semverMinor, "Fixed": semverPatch, "Security": semverPatch}

func keysOf(m changeToIncrementKind) string {
	var changes []string
	for _, i := range []incrementKind{semverMajor, semverMinor, semverPatch, semverPrerelease, semverBuild} {
		changes = append(changes, keysFor(m, i)...)
	}
	return strings.Join(changes, ", ")
}

func keysFor(m changeToIncrementKind, kind incrementKind) []string {
	var result []string
	for k, l := range m {
		if l == kind {
			result = append(result, k)
		}
	}
	sort.Strings(result)
	return result
}
