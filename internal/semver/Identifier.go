// Package semver provides helpers to work with semantic version.
// Complements "github.com/blang/semver/v4"
package semver

import "fmt"

// Identifier provides constants to reason about the parts (eg identifiers) of a semantic version
type Identifier int

const (
	Major Identifier = iota
	Minor
	Patch
	Prerelease
	Build
	endOfEnum
	startOfEnum = Major
)

func (id Identifier) String() string {
	switch id {
	case Major:
		return "major"
	case Minor:
		return "minor"
	case Patch:
		return "patch"
	case Prerelease:
		return "prerelease"
	case Build:
		return "build"
	default:
		panic(fmt.Sprintf("\"%d\" not defined", id))
	}
}

func NewIdentifier(name string) (Identifier, error) {
	switch name {
	case "major":
		return Major, nil
	case "minor":
		return Minor, nil
	case "patch":
		return Patch, nil
	case "prerelease":
		return Prerelease, nil
	case "build":
		return Build, nil
	default:
		return endOfEnum, fmt.Errorf("%q is not a valid identifier", name)
	}
}

func ForEach(f func(Identifier) error) error {
	for i := startOfEnum; i < endOfEnum; i++ {
		if err := f(i); err != nil {
			return err
		}
	}
	return nil
}
