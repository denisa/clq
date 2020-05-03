package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testHeading struct {
	name string
}

func newTestHeading(s string) Heading {
	return testHeading{name: s}
}
func (h testHeading) Name() string {
	return h.name
}
func (h testHeading) AsPath() string {
	return asPath(h.name)
}

func TestNewStackIsEmpty(t *testing.T) {
	assert := assert.New(t)

	s := NewStack()
	assert.True(s.empty(), "empty stack expected")
	assert.False(s.title(), "title not expected")
	assert.False(s.release(), "release not expected")
	assert.False(s.change(), "change not expected")
}

func TestOneLevelStackIsDocument(t *testing.T) {
	assert := assert.New(t)

	s := NewStack()
	s.push(newTestHeading("first"))
	assert.False(s.empty(), "empty stack not expected")
	assert.True(s.title(), "title expected")
	assert.False(s.release(), "release not expected")
	assert.False(s.change(), "change not expected")
}
func TestTwoLevelStackIsRelease(t *testing.T) {
	assert := assert.New(t)

	s := NewStack()
	s.push(newTestHeading("first"))
	s.push(newTestHeading("second"))
	assert.False(s.empty(), "empty stack not expected")
	assert.False(s.title(), "title not expected")
	assert.True(s.release(), "release expected")
	assert.False(s.change(), "change not expected")

}
func TestThreeLevelStackIsChange(t *testing.T) {
	assert := assert.New(t)

	s := NewStack()
	s.push(newTestHeading("first"))
	s.push(newTestHeading("second"))
	s.push(newTestHeading("third"))
	assert.False(s.empty(), "empty stack not expected")
	assert.False(s.title(), "title not expected")
	assert.False(s.release(), "release not expected")
	assert.True(s.change(), "change expected")
}

func TestPopReturnsPushedContent(t *testing.T) {
	assert := assert.New(t)

	s := NewStack()
	expected := newTestHeading("an item")
	s.push(expected)
	assert.Equal(1, s.depth())

	actual, _ := s.pop()
	assertHeadingEquals(assert, expected, actual)
	assert.Equal(0, s.depth())
}

func TestPopEmptyStackFails(t *testing.T) {
	assert := assert.New(t)

	s := NewStack()
	_, err := s.pop()
	assert.NotNil(t, err)
}
func TestPeekReturnsPushedContent(t *testing.T) {
	assert := assert.New(t)

	s := NewStack()
	expected := newTestHeading("an item")
	s.push(expected)
	assert.Equal(1, s.depth())

	actual, _ := s.Peek()
	assertHeadingEquals(assert, expected, actual)
	assert.Equal(1, s.depth())
}

func TestPeekEmptyStackFails(t *testing.T) {
	s := NewStack()
	_, err := s.Peek()
	assert.NotNil(t, err)
}

func TestStackDepthGrowsAndShrink(t *testing.T) {
	assert := assert.New(t)

	s := NewStack()
	assert.Equal(0, s.depth())
	assert.Equal("", s.AsPath())

	s.push(newTestHeading("first"))
	assert.Equal(1, s.depth())
	assert.Equal("{first}", s.AsPath())

	s.push(newTestHeading("second"))
	assert.Equal(2, s.depth())
	assert.Equal("{first}{second}", s.AsPath())

	if actual, err := s.pop(); err == nil {
		assertHeadingEquals(assert, newTestHeading("second"), actual)
	} else {
		t.Error(err)
	}
	assert.Equal(1, s.depth())
	assert.Equal("{first}", s.AsPath())

	if actual, err := s.pop(); err == nil {
		assertHeadingEquals(assert, newTestHeading("first"), actual)
	} else {
		t.Error(err)
	}

	assert.Equal(0, s.depth())
	assert.Equal("", s.AsPath())
}

func TestResetEmptyStackToZeroIsSameAsPush(t *testing.T) {
	s := NewStack()
	s.ResetTo(titleHeading, "title")
	assert.Equal(t, "{title}", s.AsPath())
}

func TestResetEmptyStackToALevelShouldFail(t *testing.T) {
	s := NewStack()
	_, err := s.ResetTo(releaseHeading, "[Unreleased]")
	assert.NotNil(t, err)
}

func TestResetFilledStackToZeroIsSameAsPushToEmptyStack(t *testing.T) {
	assert := assert.New(t)

	s := NewStack()
	if h, err := newTitle("title"); err == nil {
		s.push(h)
	}
	if h, err := newRelease("[Unreleased]"); err == nil {
		s.push(h)
	}
	if h, err := newChange("Added"); err == nil {
		s.push(h)
	}
	assert.Equal("{title}{[Unreleased]}{Added}", s.AsPath())

	s.ResetTo(titleHeading, "other title")
	assert.Equal("{other title}", s.AsPath())
}

func TestResetFilledStackToOneIsSameAsTwoPushToEmptyStack(t *testing.T) {
	assert := assert.New(t)

	s := NewStack()
	if h, err := newTitle("title"); err == nil {
		s.push(h)
	}
	if h, err := newRelease("[Unreleased]"); err == nil {
		s.push(h)
	}
	if h, err := newChange("Added"); err == nil {
		s.push(h)
	}
	assert.Equal("{title}{[Unreleased]}{Added}", s.AsPath())

	s.ResetTo(releaseHeading, "1.2.3 - 2020-04-15 [YANKED]")
	assert.Equal("{title}{1.2.3 - 2020-04-15 [YANKED]}", s.AsPath())
}
