package validator

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/blang/semver"
)

type Heading interface {
	Name() string
	AsPath() string
}

// Level 1 heading, title of the changelog
type Title struct {
	name string
}

func newTitle(s string) (Heading, error) {
	if s == "" {
		return nil, errors.New("Validation error: Title cannot stay empty")
	}
	return Title{name: s}, nil
}

func (h Title) Name() string {
	return h.name
}

func (h Title) AsPath() string {
	return asPath(h.name)
}

// Level 2 heading, release
type Release struct {
	name               string
	unreleased, yanked bool
	date               time.Time
	version            semver.Version
}

const semverPattern string = `(?P<semver>\S+?)`
const isoDatePattern string = `(?P<date>\d\d\d\d-\d\d-\d\d)`

func newRelease(s string) (Heading, error) {
	if matched, _ := regexp.MatchString(`^\[\s*Unreleased\s*\]$`, s); matched {
		return Release{name: s, unreleased: true}, nil
	}
	{
		releaseRE := regexp.MustCompile(`^\[\s*` + semverPattern + `\s*\]\s+-\s+` + isoDatePattern + `(?:\s+(?P<label>.+))?$`)
		if matches := releaseRE.FindStringSubmatch(s); matches != nil {
			groups := releaseRE.SubexpNames()
			date, err := time.Parse("2006-01-02", subexp(groups, matches, "date"))
			if err != nil {
				return nil, errors.New("Validation error: Illegal date (" + err.Error() + ") for " + s)
			}
			if matched, _ := regexp.MatchString(`[\s*YANKED\s*]`, subexp(groups, matches, "label")); matched {
				return nil, errors.New("Validation error: the version of a [YANKED] release cannot stand between [...] for " + s)
			}
			version, err := semver.Make(subexp(groups, matches, "semver"))
			if err != nil {
				return nil, errors.New("Validation error: Illegal version (" + err.Error() + ") for " + s)
			}
			return Release{name: s, date: date, version: version}, nil
		}
	}
	{
		releaseRE := regexp.MustCompile(`^` + semverPattern + `\s+-\s+` + isoDatePattern + `\s+\[\s*YANKED\s*]?$`)
		if matches := releaseRE.FindStringSubmatch(s); matches != nil {
			groups := releaseRE.SubexpNames()
			date, err := time.Parse("2006-01-02", subexp(groups, matches, "date"))
			if err != nil {
				return nil, errors.New("Validation error: Illegal date (" + err.Error() + ") for " + s)
			}
			version, err := semver.Make(subexp(groups, matches, "semver"))
			if err != nil {
				return nil, errors.New("Validation error: Illegal version (" + err.Error() + ") for " + s)
			}
			return Release{name: s, date: date, version: version, yanked: true}, nil
		}
	}
	return nil, errors.New("Validation error: Unknown release header for " + s)
}

func (h Release) Name() string {
	return h.name
}

func (h Release) AsPath() string {
	return asPath(h.name)
}

func (r *Release) isRelease() bool {
	return r.date != time.Time{} && !r.unreleased && len(r.version.Pre) == 0 && len(r.version.Build) == 0
}

const (
	semverMajor = iota
	semverMinor
	semverPatch
	semverPrerelease
	semverBuild
)

// Level 3 heading, change groups
type Change struct {
	name string
}

func newChange(s string) (Heading, error) {
	for val := range changeKind {
		if matched, _ := regexp.MatchString(`^`+val+`$`, s); matched {
			return Change{name: s}, nil
		}
	}

	return nil, errors.New("Validation error: Unknown change headings '" + s + "' is not one of [" + keysOf(changeKind) + "]")
}

func (h Change) Name() string {
	return h.name
}
func (h Change) AsPath() string {
	return asPath(h.name)
}

var changeKind = map[string]int{"Added": semverMajor, "Removed": semverMajor, "Changed": semverMinor, "Deprecated": semverMinor, "Fixed": semverPatch, "Security": semverPatch}

func keysOf(m map[string]int) string {
	var keys strings.Builder
	for k := range m {
		if keys.Len() > 0 {
			keys.WriteString(", ")
		}
		keys.WriteString("\"")
		keys.WriteString(k)
		keys.WriteString("\"")
	}
	return keys.String()
}

func keysFor(m map[string]int, level int) []string {
	var result []string
	for k, l := range m {
		if l == level {
			result = append(result, k)
		}
	}
	return result
}

func nextRelease(c changeMap, v semver.Version) semver.Version {
	for _, k := range keysFor(changeKind, semverMajor) {
		if c[k] {
			return semver.Version{Major: v.Major + 1, Minor: 0, Patch: 0}
		}
	}
	for _, k := range keysFor(changeKind, semverMinor) {
		if c[k] {
			return semver.Version{Major: v.Major, Minor: v.Minor + 1, Patch: 0}
		}
	}
	for _, k := range keysFor(changeKind, semverPatch) {
		if c[k] {
			return semver.Version{Major: v.Major, Minor: v.Minor, Patch: v.Patch + 1}
		}
	}
	return v
}

func asPath(name string) string {
	return "{" + name + "}"
}

func subexp(groups []string, matches []string, subexp string) string {
	for index, name := range groups {
		if name == subexp {
			return matches[index]
		}
	}

	return ""
}
