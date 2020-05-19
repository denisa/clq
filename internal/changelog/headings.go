package changelog

import (
	"fmt"
	"sort"
	"strings"
)

type HeadingKind int

const (
	TitleHeading HeadingKind = iota
	ReleaseHeading
	ChangeHeading
)

type Heading interface {
	Name() string
	String() string
}

func asPath(name string) string {
	return "{" + name + "}"
}

func NewHeading(level HeadingKind, name string) (Heading, error) {
	switch level {
	case TitleHeading:
		return newChangelog(name)
	case ReleaseHeading:
		return newRelease(name)
	case ChangeHeading:
		return newChange(name)
	}
	return nil, fmt.Errorf("Unknown heading level %v", level)
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

func keysFor(m changeToIncrementKind, level incrementKind) []string {
	var result []string
	for k, l := range m {
		if l == level {
			result = append(result, k)
		}
	}
	sort.Strings(result)
	return result
}
