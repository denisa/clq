package validator

import (
	"errors"
	"fmt"
	"strings"
)

type stack struct {
	s []string
}

func NewStack() stack {
	return stack{make([]string, 0)}
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

func (s *stack) push(v string) {
	s.s = append(s.s, v)
}

func (s *stack) pop() (string, error) {
	l := s.depth()
	if l == 0 {
		return "", errors.New("Empty stack")
	}
	res := s.s[l-1]
	s.s = s.s[:l-1]
	return res, nil
}

func (s *stack) resetTo(depth int, name string) error {
	if depth > s.depth() {
		return fmt.Errorf("Attempting to reset to %d for a stack of depth %d", depth, s.depth())
	}
	for i := s.depth(); i > depth; i-- {
		s.pop()
	}
	s.push(name)
	return nil
}

func (s *stack) asPath() string {
	var path strings.Builder
	for _, val := range s.s {
		path.WriteString("{")
		path.WriteString(val)
		path.WriteString("}")
	}
	return path.String()
}
