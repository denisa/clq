package changelog

import (
	"fmt"
	"regexp"
	"time"

	"github.com/blang/semver/v4"
	semverConstants "github.com/denisa/clq/internal/semver"
)

type Release struct {
	heading
	label   string
	yanked  bool
	date    time.Time
	version semver.Version
}

const semverPattern string = `(?P<semver>\S+)`
const isoDatePattern string = `(?P<date>\S+)`

func (h HeadingsFactory) newRelease(title string) (Heading, error) {
	if matched, _ := regexp.MatchString(`^\[\s*Unreleased\s*]$`, title); matched {
		return Release{heading: heading{title: title, kind: ReleaseHeading}}, nil
	}
	{
		releaseRE := regexp.MustCompile(`^\[\s*` + semverPattern + `\s*\]\s+-\s+` + isoDatePattern + `(?:\s+(?P<label>.+))?$`)
		if matches := releaseRE.FindStringSubmatch(title); matches != nil {
			groups := releaseRE.SubexpNames()
			date, err := time.Parse("2006-01-02", subexp(groups, matches, "date"))
			if err != nil {
				return nil, fmt.Errorf("validation error: Illegal date (%v) for %v", err, title)
			}
			label := subexp(groups, matches, "label")
			if matched, _ := regexp.MatchString(`\[\s*YANKED\s*]`, label); matched {
				return nil, fmt.Errorf("validation error: the version of a [YANKED] release cannot stand between [...] for %v", title)
			}
			version, err := semver.Make(subexp(groups, matches, "semver"))
			if err != nil {
				return nil, fmt.Errorf("validation error: Illegal version (%v) for %v", err, title)
			}
			return Release{heading: heading{title: title, kind: ReleaseHeading}, date: date, version: version, label: label}, nil
		}
	}
	{
		releaseRE := regexp.MustCompile(`^` + semverPattern + `\s+-\s+` + isoDatePattern + `\s+\[\s*YANKED\s*]?$`)
		if matches := releaseRE.FindStringSubmatch(title); matches != nil {
			groups := releaseRE.SubexpNames()
			date, err := time.Parse("2006-01-02", subexp(groups, matches, "date"))
			if err != nil {
				return nil, fmt.Errorf("validation error: Illegal date (%v) for %v", err, title)
			}
			version, err := semver.Make(subexp(groups, matches, "semver"))
			if err != nil {
				return nil, fmt.Errorf("validation error: Illegal version (%v) for %v", err, title)
			}
			return Release{heading: heading{title: title, kind: ReleaseHeading}, date: date, version: version, yanked: true}, nil
		}
	}
	return nil, fmt.Errorf("validation error: Unknown release header for %q", title)
}

func (h Release) DisplayTitle() string {
	return h.Title()
}
func (h Release) Title() string {
	return h.title
}

func (h Release) Kind() HeadingKind {
	return h.kind
}

func (h Release) String() string {
	return asPath(h.title)
}

// Date returns the release date if this has been released, an empty string otherwise.
func (h Release) Date() string {
	if h.HasBeenReleased() {
		return h.date.Format("2006-01-02")
	}
	return ""
}

// Label returns the optional label of this release.
func (h Release) Label() string {
	return h.label
}

// Version returns the release version if this has been released, an empty string otherwise.
func (h Release) Version() string {
	if h.HasBeenReleased() {
		return h.version.String()
	}
	return ""
}

// IsPrerelease returns true if this release is a pre-release without build number component.
func (h Release) IsPrerelease() bool {
	return h.HasBeenReleased() && len(h.version.Pre) > 0 && len(h.version.Build) == 0
}

// IsRelease returns true if this release is a release without pre-release or build number component.
func (h Release) IsRelease() bool {
	return h.HasBeenReleased() && len(h.version.Pre) == 0 && len(h.version.Build) == 0
}

// HasBeenYanked returns true if this release has been yanked.
func (h Release) HasBeenYanked() bool {
	return h.yanked
}

// HasBeenReleased returns true if this has ever been released
func (h Release) HasBeenReleased() bool {
	return !h.date.IsZero()
}

// ReleaseIs returns true this release is equal to the otherRelease
func (h Release) ReleaseIs(otherRelease semver.Version) bool {
	return h.version.EQ(otherRelease)
}

// NextRelease computes what the next version number should be given a set of changes.
func (h Release) NextRelease(semverIdentifier semverConstants.Identifier) semver.Version {
	switch semverIdentifier {
	case semverConstants.Major:
		return semver.Version{Major: h.version.Major + 1, Minor: 0, Patch: 0}
	case semverConstants.Minor:
		return semver.Version{Major: h.version.Major, Minor: h.version.Minor + 1, Patch: 0}
	case semverConstants.Patch:
		return semver.Version{Major: h.version.Major, Minor: h.version.Minor, Patch: h.version.Patch + 1}
	default:
		return h.version
	}
}

// IsNewerThan returns an error if this release if not newer than another release.
// A release is newer if it has the same or a more recent date and one its
// version has been incremented.
func (h Release) IsNewerThan(other Release) error {
	if h.date.Before(other.date) {
		return fmt.Errorf("validation error: release %q should be older than %q", other.Title(), h.Title())
	}
	if h.version.LTE(other.version) {
		return fmt.Errorf("validation error: release %q should sort before %q", other.Title(), h.Title())
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
