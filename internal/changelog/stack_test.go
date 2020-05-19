package changelog

import (
	"testing"

	"github.com/stretchr/testify/require"
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
	require := require.New(t)

	s := NewStack()
	require.True(s.empty(), "empty stack expected")
	require.False(s.Title(), "title not expected")
	require.False(s.Release(), "release not expected")
	require.False(s.Change(), "change not expected")
}

func TestOneLevelStackIsDocument(t *testing.T) {
	require := require.New(t)

	s := NewStack()
	s.push(newTestHeading("first"))
	require.False(s.empty(), "empty stack not expected")
	require.True(s.Title(), "title expected")
	require.False(s.Release(), "release not expected")
	require.False(s.Change(), "change not expected")
}
func TestTwoLevelStackIsRelease(t *testing.T) {
	require := require.New(t)

	s := NewStack()
	s.push(newTestHeading("first"))
	s.push(newTestHeading("second"))
	require.False(s.empty(), "empty stack not expected")
	require.False(s.Title(), "title not expected")
	require.True(s.Release(), "release expected")
	require.False(s.Change(), "change not expected")

}
func TestThreeLevelStackIsChange(t *testing.T) {
	require := require.New(t)

	s := NewStack()
	s.push(newTestHeading("first"))
	s.push(newTestHeading("second"))
	s.push(newTestHeading("third"))
	require.False(s.empty(), "empty stack not expected")
	require.False(s.Title(), "title not expected")
	require.False(s.Release(), "release not expected")
	require.True(s.Change(), "change expected")
}

func TestPopReturnsPushedContent(t *testing.T) {
	require := require.New(t)

	s := NewStack()
	expected := newTestHeading("an item")
	s.push(expected)
	require.Equal(1, s.depth())

	actual, _ := s.pop()
	requireHeadingEquals(require, expected, actual)
	require.Equal(0, s.depth())
}

func TestPopEmptyStackFails(t *testing.T) {
	require := require.New(t)

	s := NewStack()
	_, err := s.pop()
	require.NotNil(t, err)
}

func TestStackDepthGrowsAndShrink(t *testing.T) {
	require := require.New(t)

	s := NewStack()
	require.Equal(0, s.depth())
	require.Equal("", s.AsPath())

	s.push(newTestHeading("first"))
	require.Equal(1, s.depth())
	require.Equal("{first}", s.AsPath())

	s.push(newTestHeading("second"))
	require.Equal(2, s.depth())
	require.Equal("{first}{second}", s.AsPath())

	if actual, err := s.pop(); err == nil {
		requireHeadingEquals(require, newTestHeading("second"), actual)
	} else {
		t.Error(err)
	}
	require.Equal(1, s.depth())
	require.Equal("{first}", s.AsPath())

	if actual, err := s.pop(); err == nil {
		requireHeadingEquals(require, newTestHeading("first"), actual)
	} else {
		t.Error(err)
	}

	require.Equal(0, s.depth())
	require.Equal("", s.AsPath())
}

func TestResetEmptyStackToZeroIsSameAsPush(t *testing.T) {
	s := NewStack()
	s.ResetTo(TitleHeading, "title")
	require.Equal(t, "{title}", s.AsPath())
}

func TestResetEmptyStackToALevelShouldFail(t *testing.T) {
	s := NewStack()
	_, err := s.ResetTo(ReleaseHeading, "[Unreleased]")
	require.NotNil(t, err)
}

func TestResetFilledStackToZeroIsSameAsPushToEmptyStack(t *testing.T) {
	require := require.New(t)

	s := NewStack()
	if h, err := newChangelog("title"); err == nil {
		s.push(h)
	}
	if h, err := newRelease("[Unreleased]"); err == nil {
		s.push(h)
	}
	if h, err := newChange("Added"); err == nil {
		s.push(h)
	}
	require.Equal("{title}{[Unreleased]}{Added}", s.AsPath())

	s.ResetTo(TitleHeading, "other title")
	require.Equal("{other title}", s.AsPath())
}

func TestResetFilledStackToOneIsSameAsTwoPushToEmptyStack(t *testing.T) {
	require := require.New(t)

	s := NewStack()
	if h, err := newChangelog("title"); err == nil {
		s.push(h)
	}
	if h, err := newRelease("[Unreleased]"); err == nil {
		s.push(h)
	}
	if h, err := newChange("Added"); err == nil {
		s.push(h)
	}
	require.Equal("{title}{[Unreleased]}{Added}", s.AsPath())

	s.ResetTo(ReleaseHeading, "1.2.3 - 2020-04-15 [YANKED]")
	require.Equal("{title}{1.2.3 - 2020-04-15 [YANKED]}", s.AsPath())
}
