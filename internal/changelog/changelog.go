// Package changelog keeps track of the current state during the traversal of a changelog document.
// Things are named according to entries in the grammar.
package changelog

import (
	"fmt"
	"strings"
)

// Changelog is the current state of the changelog during its traversal.
type Changelog struct {
	headingsFactory HeadingsFactory
	headings        []Heading
	listeners       []Listener
}

// NewChangelog creates a new empty changelog.
func NewChangelog(headingsFactory HeadingsFactory) *Changelog {
	return &Changelog{headingsFactory: headingsFactory}
}

// Listener registers one or more listeners to this changelog.
// Listeners are notified as sections are entered and exited.
func (c *Changelog) Listener(l ...Listener) {
	c.listeners = append(c.listeners, l...)
}

// Introduction returns true if we’re currently visiting the Introduction section.
func (c *Changelog) Introduction() bool {
	return len(c.headings) > 0 && c.active() == IntroductionHeading
}

// Release returns true if we’re currently visiting one of the Release sections.
func (c *Changelog) Release() bool {
	return len(c.headings) > 0 && c.active() == ReleaseHeading
}

// Change returns true if we’re currently visiting one of the Change sections.
func (c *Changelog) Change() bool {
	return len(c.headings) > 0 && (c.active() == ChangeHeading || c.active() == ChangeDescription)
}

// active returns the kind of the top-most heading
func (c *Changelog) active() HeadingKind {
	return c.headings[len(c.headings)-1].Kind()
}

// Close calls the registered Exit listeners for all the un-closed headings.
func (c *Changelog) Close() {
	for i := HeadingKind(len(c.headings)) - 1; i > -1; i-- {
		for _, l := range c.listeners {
			l.Exit(c.headings[i])
		}
	}
}

// Section sets the state to a new section kind with the given title.
// Section can go down one-level, for example from Release to Change, or up any number of levels.
// Section creates and returns the section’s Heading
func (c *Changelog) Section(kind HeadingKind, title string) (Heading, error) {
	if kind > HeadingKind(len(c.headings)) {
		return nil, fmt.Errorf("attempting to roll-back a changelog at %v to %v", len(c.headings), kind)
	}

	h, err := c.headingsFactory.NewHeading(kind, title)
	if err != nil {
		return nil, err
	}

	for i := HeadingKind(len(c.headings)) - 1; i >= kind; i-- {
		for _, l := range c.listeners {
			l.Exit(c.headings[i])
		}
	}

	c.headings = append(c.headings[:kind], h)
	for _, l := range c.listeners {
		l.Enter(h)
	}
	return h, nil
}

func (c *Changelog) String() string {
	var path strings.Builder
	for _, heading := range c.headings {
		path.WriteString(heading.String())
	}
	return path.String()
}
