package changelog

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/denisa/clq/internal/semver"
)

type config struct {
	semver semver.Identifier
	emoji string
}
type changeKindToConfig map[string]config

// Collects information about the supported change headers
type ChangeKind struct {
	changes changeKindToConfig
}

// Loads a new ChangeKind from a file
func NewChangeKind(fileName string) (*ChangeKind, error) {
	if fileName == "" {
		return newChangeKind(), nil
	}

	file, e := ioutil.ReadFile(fileName)
	if e != nil {
		return nil, e
	}

	c := &ChangeKind{changes: make(changeKindToConfig)}
	if err := json.Unmarshal(file, c); err != nil {
		return nil, err
	}
	return c, nil
}

func newChangeKind() *ChangeKind {
	return &ChangeKind{changes: changeKindToConfig{"Added": {semver.Major, ""}, "Removed": {semver.Major, ""}, "Changed": {semver.Minor, ""}, "Deprecated": {semver.Minor, ""}, "Fixed": {semver.Patch, ""}, "Security": {semver.Patch, ""}}}
}

// Returns an error if the given title is not a recognized change kind.
func (m ChangeKind) IsSupported(title string) error {
	if _, ok := m.changes[title]; ok {
		return nil
	}
	return fmt.Errorf("Validation error: Unknown change heading %q is not one of [%v]", title, m.keysOf())
}

// Returns the increment kind to apply for a set of change kinds and the reason for it.
func (m ChangeKind) IncrementFor(c ChangeMap) (semver.Identifier, string) {
	increment := semver.Build
	trigger := ""
	for k := range c {
		if val, ok := m.changes[k]; ok && val.semver < increment {
			increment = val.semver
			trigger = k
		}
	}
	return increment, trigger
}

func (m ChangeKind) add(name string, increment semver.Identifier, emoji string) error {
	if strings.TrimSpace(name) == "" {
		return  fmt.Errorf("Validation error: \"name\" is blank")
	}

	m.changes[name] = config{ semver: increment, emoji: emoji }
	return nil
}

func (m ChangeKind) keysOf() string {
	var result []string
	semver.ForEach(func(i semver.Identifier) error {
		result = append(result, m.keysFor(i)...)
		return nil
	})
	sort.Strings(result)
	return strings.Join(result, ", ")
}

func (m ChangeKind) keysFor(kind semver.Identifier) []string {
	var result []string
	for k, l := range m.changes {
		if l.semver == kind {
			result = append(result, k)
		}
	}
	return result
}

func (m ChangeKind) emojiFor(title string) (string, error) {
	if c, ok := m.changes[title]; ok {
		return c.emoji, nil
	}
	return "", fmt.Errorf("Unknown change heading %q is not one of [%v]", title, m.keysOf())
}
