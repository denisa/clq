package changelog

import (
	"fmt"
	"strings"
)

type Changelog struct {
	s []Heading
}

func NewChangelog() Changelog {
	return Changelog{make([]Heading, 0)}
}

func (s *Changelog) Title() bool {
	return len(s.s) == 1
}

func (s *Changelog) Release() bool {
	return len(s.s) == 2
}

func (s *Changelog) Change() bool {
	return len(s.s) == 3
}

func (s *Changelog) ResetTo(kind HeadingKind, name string) (Heading, error) {
	if kind > HeadingKind(len(s.s)) {
		return nil, fmt.Errorf("Attempting to roll-back a changelog at %v to %v", len(s.s), kind)
	}

	h, err := NewHeading(kind, name)
	if err != nil {
		return nil, err
	}

	s.s = append(s.s[:kind], h)
	return h, nil
}

func (s Changelog) String() string {
	var path strings.Builder
	for _, heading := range s.s {
		path.WriteString(heading.String())
	}
	return path.String()
}
