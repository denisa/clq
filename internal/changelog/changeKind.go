package changelog

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/denisa/clq/internal/semver"
)

type config struct {
	semver semver.Identifier
	emoji  string
}
type changeKindToConfig map[string]config

// ChangeKind collects information about the supported change headers
type ChangeKind struct {
	changes changeKindToConfig
}

// NewChangeKind loads a new ChangeKind from a file
func NewChangeKind(fileName string) (*ChangeKind, error) {
	if fileName == "" {
		return &ChangeKind{changes: changeKindToConfig{"Added": {semver.Major, ""}, "Removed": {semver.Major, ""}, "Changed": {semver.Minor, ""}, "Deprecated": {semver.Minor, ""}, "Fixed": {semver.Patch, ""}, "Security": {semver.Patch, ""}}}, nil
	}

	file, e := os.ReadFile(fileName)
	if e != nil {
		return nil, e
	}

	c := &ChangeKind{changes: make(changeKindToConfig)}
	if err := json.Unmarshal(file, c); err != nil {
		return nil, err
	}
	return c, nil
}

// IncrementFor returns the increment kind to apply for a set of change kinds and the reason for it.
func (ck ChangeKind) IncrementFor(c ChangeMap) (semver.Identifier, string) {
	increment := semver.Build
	trigger := ""
	for k := range c {
		if val, ok := ck.changes[k]; ok && val.semver < increment {
			increment = val.semver
			trigger = k
		}
	}
	return increment, trigger
}

func (ck ChangeKind) add(name string, increment semver.Identifier, emoji string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("validation error: \"name\" is blank")
	}

	ck.changes[name] = config{semver: increment, emoji: emoji}
	return nil
}

func (ck ChangeKind) keysOf() string {
	var result []string
	_ = semver.ForEach(func(i semver.Identifier) error {
		result = append(result, ck.keysFor(i)...)
		return nil
	})
	sort.Strings(result)
	return strings.Join(result, ", ")
}

func (ck ChangeKind) keysFor(kind semver.Identifier) []string {
	var result []string
	for k, l := range ck.changes {
		if l.semver == kind {
			result = append(result, k)
		}
	}
	return result
}

func (ck ChangeKind) emojiFor(title string) (string, error) {
	if c, ok := ck.changes[title]; ok {
		return c.emoji, nil
	}
	return "", fmt.Errorf("validation error: Unknown change heading %q is not one of [%v]", title, ck.keysOf())
}
