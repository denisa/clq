package changelog

import (
	"fmt"
	"regexp"
	"time"

	"github.com/blang/semver"
)

// Level 2, release
type ChangeMap map[string]bool

type Release struct {
	name, label        string
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
				return nil, fmt.Errorf("Validation error: Illegal date (%v) for %v", err, s)
			}
			label := subexp(groups, matches, "label")
			if matched, _ := regexp.MatchString(`[\s*YANKED\s*]`, label); matched {
				return nil, fmt.Errorf("Validation error: the version of a [YANKED] release cannot stand between [...] for %v", s)
			}
			version, err := semver.Make(subexp(groups, matches, "semver"))
			if err != nil {
				return nil, fmt.Errorf("Validation error: Illegal version (%v) for %v", err, s)
			}
			return Release{name: s, date: date, version: version, label: label}, nil
		}
	}
	{
		releaseRE := regexp.MustCompile(`^` + semverPattern + `\s+-\s+` + isoDatePattern + `\s+\[\s*YANKED\s*]?$`)
		if matches := releaseRE.FindStringSubmatch(s); matches != nil {
			groups := releaseRE.SubexpNames()
			date, err := time.Parse("2006-01-02", subexp(groups, matches, "date"))
			if err != nil {
				return nil, fmt.Errorf("Validation error: Illegal date (%v) for %v", err, s)
			}
			version, err := semver.Make(subexp(groups, matches, "semver"))
			if err != nil {
				return nil, fmt.Errorf("Validation error: Illegal version (%v) for %v", err, s)
			}
			return Release{name: s, date: date, version: version, yanked: true}, nil
		}
	}
	return nil, fmt.Errorf("Validation error: Unknown release header for %q", s)
}

func (h Release) Name() string {
	return h.name
}

func (h Release) String() string {
	return asPath(h.name)
}

func (h Release) Date() string {
	return h.date.Format("2006-01-02")
}

func (h Release) Label() string {
	return h.label
}

func (h Release) Version() string {
	return h.version.String()
}

func (h Release) IsRelease() bool {
	return h.HasBeenReleased() && !h.unreleased && len(h.version.Pre) == 0 && len(h.version.Build) == 0
}

func (h Release) Unreleased() bool {
	return h.unreleased
}

func (h Release) Yanked() bool {
	return h.yanked
}

func (h Release) HasBeenReleased() bool {
	return h.date != time.Time{}
}

func (h Release) HasRelease(nextRelease semver.Version) bool {
	return h.version.EQ(nextRelease)
}

func (h Release) NextRelease(c ChangeMap) semver.Version {
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

func (h Release) SortsBefore(other Release) error {
	if h.date.Before(other.date) {
		return fmt.Errorf("Validation error: release %q should be older than %q", other.Name(), h.Name())
	}
	if h.version.LTE(other.version) {
		return fmt.Errorf("Validation error: release %q should sort before %q", other.Name(), h.Name())
	}
	return nil
}

func subexp(groups []string, matches []string, subexp string) string {
	for index, name := range groups {
		if name == subexp {
			return matches[index]
		}
	}

	return ""
}
