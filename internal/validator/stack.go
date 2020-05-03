package validator

import (
	"errors"
	"fmt"
	"strings"
)

const (
	titleHeading = iota
	releaseHeading
	changeHeading
)

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
func (s *stack) ResetTo(depth int, name string) (Heading, error) {
	if depth > s.depth() {
		return nil, fmt.Errorf("Attempting to reset to %d for a stack of depth %d", depth, s.depth())
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
			return nil, err
		}
	}

	s.s = s.s[:depth]
	s.push(h)
	return h, nil
}

func (s *stack) AsPath() string {
	var path strings.Builder
	for _, heading := range s.s {
		path.WriteString(heading.AsPath())
	}
	return path.String()
}
