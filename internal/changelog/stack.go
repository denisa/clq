package changelog

import (
	"fmt"
	"strings"
)

type Stack struct {
	s []Heading
}

func NewStack() Stack {
	return Stack{make([]Heading, 0)}
}

func (s *Stack) empty() bool {
	return s.depth() == 0
}

func (s *Stack) Title() bool {
	return s.depth() == 1
}

func (s *Stack) Release() bool {
	return s.depth() == 2
}

func (s *Stack) Change() bool {
	return s.depth() == 3
}

func (s *Stack) depth() int {
	return len(s.s)
}

func (s *Stack) push(v Heading) {
	s.s = append(s.s, v)
}

func (s *Stack) pop() (Heading, error) {
	l := s.depth()
	if l == 0 {
		return nil, fmt.Errorf("Empty stack")
	}
	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}

func (s *Stack) Peek() (Heading, error) {
	l := s.depth()
	if l == 0 {
		return nil, fmt.Errorf("Empty stack")
	}
	return s.s[l-1], nil
}

func (s *Stack) ResetTo(depth HeadingKind, name string) (Heading, error) {
	if depth > HeadingKind(s.depth()) {
		return nil, fmt.Errorf("Attempting to reset to %d for a stack of depth %d", depth, s.depth())
	}

	if h, err := NewHeading(depth, name); err != nil {
		return nil, err
	} else {
		s.s = s.s[:depth]
		s.push(h)
		return h, nil
	}
}

func (s *Stack) AsPath() string {
	var path strings.Builder
	for _, heading := range s.s {
		path.WriteString(heading.AsPath())
	}
	return path.String()
}
