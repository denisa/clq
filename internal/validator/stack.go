package validator

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/blang/semver"
)

const (
	titleHeading = iota
	releaseHeading
	changeHeading
)

const semverPattern string = `(?P<semver>\S+?)`
const isoDatePattern string = `(?P<date>\d\d\d\d-\d\d-\d\d)`

type Heading interface {
	Name() string
	AsPath() string
}

func asPath(name string) string {
	return "{" + name + "}"
}

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

type Release struct {
	name               string
	unreleased, yanked bool
	date               time.Time
	version            semver.Version
}

func newRelease(s string) (Heading, error) {
	if matched, _ := regexp.MatchString(`^\[\s*Unreleased\s*\]$`, s); matched {
		return Release{name: s, unreleased: true}, nil
	}
	{
		releaseRE := regexp.MustCompile(`^\[\s*` + semverPattern + `\s*\]\s+-\s+` + isoDatePattern + `(?:\s+(?P<label>.+))?$`)
		if matches := releaseRE.FindStringSubmatch(s); matches != nil {
			date, err := time.Parse("2006-01-02", subexp(releaseRE, matches, "date"))
			if err != nil {
				return nil, errors.New("Validation error: Illegal date (" + err.Error() + ") for " + s)
			}
			if matched, _ := regexp.MatchString(`[\s*YANKED\s*]`, subexp(releaseRE, matches, "label")); matched {
				return nil, errors.New("Validation error: the version of a [YANKED] release cannot stand between [...] for " + s)
			}
			version, err := semver.Make(subexp(releaseRE, matches, "semver"))
			if err != nil {
				return nil, errors.New("Validation error: Illegal version (" + err.Error() + ") for " + s)
			}
			return Release{name: s, date: date, version: version}, nil
		}
	}
	{
		yankedReleaseRE := regexp.MustCompile(`^` + semverPattern + `\s+-\s+` + isoDatePattern + `\s+\[\s*YANKED\s*]?$`)
		if matches := yankedReleaseRE.FindStringSubmatch(s); matches != nil {
			date, err := time.Parse("2006-01-02", subexp(yankedReleaseRE, matches, "date"))
			if err != nil {
				return nil, errors.New("Validation error: Illegal date (" + err.Error() + ") for " + s)
			}
			version, err := semver.Make(subexp(yankedReleaseRE, matches, "semver"))
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

type Change struct {
	name string
}

func newChange(s string) (Heading, error) {
	var changes = []string{"Added", "Removed", "Changed", "Deprecated", "Fixed", "Security"}
	for _, val := range changes {
		if matched, _ := regexp.MatchString(`^`+val+`$`, s); matched {
			return Change{name: s}, nil
		}
	}
	return nil, errors.New("Validation error: Unknown change headings '" + s + "' not supported")
}

func (h Change) Name() string {
	return h.name
}
func (h Change) AsPath() string {
	return asPath(h.name)
}

type stack struct {
	s []Heading
}

func NewStack() stack {
	return stack{make([]Heading, 0)}
}

func (s *stack) empty() bool {
	return s.depth() == 0
}
func (s *stack) title() bool {
	return s.depth() == 1
}

func (s *stack) release() bool {
	return s.depth() == 2
}

func (s *stack) change() bool {
	return s.depth() == 3
}

func (s *stack) depth() int {
	return len(s.s)
}

func (s *stack) push(v Heading) {
	s.s = append(s.s, v)
}

func (s *stack) pop() (Heading, error) {
	l := s.depth()
	if l == 0 {
		return nil, errors.New("Empty stack")
	}
	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}

func (s *stack) Peek() (Heading, error) {
	l := s.depth()
	if l == 0 {
		return nil, errors.New("Empty stack")
	}
	return s.s[l-1], nil
}
func (s *stack) ResetTo(depth int, name string) error {
	if depth > s.depth() {
		return fmt.Errorf("Attempting to reset to %d for a stack of depth %d", depth, s.depth())
	}
	var h Heading
	{
		var err error
		switch depth {
		case titleHeading:
			h, err = newTitle(name)
		case releaseHeading:
			h, err = newRelease(name)
		case changeHeading:
			h, err = newChange(name)
		}
		if err != nil {
			return err
		}
	}

	for i := s.depth(); i > depth; i-- {
		s.pop()
	}
	s.push(h)
	return nil
}

func (s *stack) AsPath() string {
	var path strings.Builder
	for _, heading := range s.s {
		path.WriteString(heading.AsPath())
	}
	return path.String()
}

func subexp(exp *regexp.Regexp, matches []string, subexp string) string {
	for index, name := range exp.SubexpNames() {
		if index >= len(matches) {
			continue
		}

		if name == subexp {
			return matches[index]
		}
	}

	return ""
}
