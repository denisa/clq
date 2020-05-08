package changelog

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/blang/semver"
)

const (
	TitleHeading = iota
	ReleaseHeading
	ChangeHeading
)

type Heading interface {
	Name() string
	AsPath() string
}

func NewHeading(level int, name string) (Heading, error) {
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

// Level 1, changelog
type Changelog struct {
	title string
}

func newChangelog(s string) (Heading, error) {
	if s == "" {
		return nil, errors.New("Validation error: Changelog cannot stay empty")
	}
	return Changelog{title: s}, nil
}

func (h Changelog) Name() string {
	return h.title
}

func (h Changelog) AsPath() string {
	return asPath(h.title)
}

// Level 2, release
type ChangeMap map[string]bool

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

func (h *Release) IsRelease() bool {
	return h.HasBeenReleased() && !h.unreleased && len(h.version.Pre) == 0 && len(h.version.Build) == 0
}

func (h *Release) Unreleased() bool {
	return h.unreleased
}

func (h *Release) Yanked() bool {
	return h.yanked
}

func (h *Release) HasBeenReleased() bool {
	return h.date != time.Time{}
}

func (h *Release) HasRelease(nextRelease semver.Version) bool {
	return h.version.EQ(nextRelease)
}

func (h *Release) NextRelease(c ChangeMap) semver.Version {
	for _, k := range keysFor(changeKind, semverMajor) {
		if c[k] {
			return semver.Version{Major: h.version.Major + 1, Minor: 0, Patch: 0}
		}
	}
	for _, k := range keysFor(changeKind, semverMinor) {
		if c[k] {
			return semver.Version{Major: h.version.Major, Minor: h.version.Minor + 1, Patch: 0}
		}
	}
	for _, k := range keysFor(changeKind, semverPatch) {
		if c[k] {
			return semver.Version{Major: h.version.Major, Minor: h.version.Minor, Patch: h.version.Patch + 1}
		}
	}
	return h.version
}

func (h *Release) SortsBefore(other Release) error {
	if h.date.Before(other.date) {
		return errors.New("Validation error: release '" + other.Name() + "' should be older than '" + h.Name() + "'")
	}
	if h.version.LTE(other.version) {
		return errors.New("Validation error: release '" + other.Name() + "' should be older than '" + h.Name() + "'")
	}
	return nil
}

const (
	semverMajor = iota
	semverMinor
	semverPatch
	semverPrerelease
	semverBuild
)

// Level 3, change groups
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
	sep := ""
	for k := range m {
		keys.WriteString(sep)
		keys.WriteString("\"")
		keys.WriteString(k)
		keys.WriteString("\"")
		sep = ", "
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
