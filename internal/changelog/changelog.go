// The changelog package keeps track of the current state during the traversal of a changelog document.
// Things are named according to entries in the grammar.
package changelog

import (
	"fmt"
	"strings"
)

// Changelog is the current state of the changelog during its traversal.
type Changelog struct {
	s []Heading
}

// NewChangelog creates a new empty changelog.
func NewChangelog() Changelog {
	return Changelog{}
}

// Introduction returns true if we’re currently visiting the Introduction section.
func (c *Changelog) Introduction() bool {
	return len(c.s) == 1
}

// Release returns true if we’re currently visiting one of the Release sections.
func (c *Changelog) Release() bool {
	return len(c.s) == 2
}

// Change returns true if we’re currently visiting one of the Change sections.
func (c *Changelog) Change() bool {
	return len(c.s) == 3
}

// Section sets the state to a new section kind with the given title.
// Section can go down one-level, for example from Release to Change, or up any number of levels.
// Section creates and returns the section’s Heading
func (c *Changelog) Section(kind HeadingKind, title string) (Heading, error) {
	if kind > HeadingKind(len(c.s)) {
		return nil, fmt.Errorf("Attempting to roll-back a changelog at %v to %v", len(c.s), kind)
	}

	h, err := NewHeading(kind, title)
	if err != nil {
		return nil, err
	}

	c.s = append(c.s[:kind], h)
	return h, nil
}

func (c Changelog) String() string {
	var path strings.Builder
	for _, heading := range c.s {
		path.WriteString(heading.String())
	}
	return path.String()
}
