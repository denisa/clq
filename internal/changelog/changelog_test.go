package changelog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewChangelogIsEmpty(t *testing.T) {
	require := require.New(t)

	s := NewChangelog()

	require.False(s.Title(), "title not expected")
	require.False(s.Release(), "release not expected")
	require.False(s.Change(), "change not expected")
}

func TestOneLevelChangelogIsDocument(t *testing.T) {
	require := require.New(t)

	s := NewChangelog()
	s.ResetTo(IntroductionHeading, "title")

	require.True(s.Title(), "title expected")
	require.False(s.Release(), "release not expected")
	require.False(s.Change(), "change not expected")
}

func TestTwoLevelChangelogIsRelease(t *testing.T) {
	require := require.New(t)

	s := NewChangelog()
	s.ResetTo(IntroductionHeading, "title")
	s.ResetTo(ReleaseHeading, "[Unreleased]")

	require.False(s.Title(), "title not expected")
	require.True(s.Release(), "release expected")
	require.False(s.Change(), "change not expected")

}

func TestThreeLevelChangelogIsChange(t *testing.T) {
	require := require.New(t)

	s := NewChangelog()
	s.ResetTo(IntroductionHeading, "title")
	s.ResetTo(ReleaseHeading, "[Unreleased]")
	s.ResetTo(ChangeHeading, "Added")

	require.False(s.Title(), "title not expected")
	require.False(s.Release(), "release not expected")
	require.True(s.Change(), "change expected")
}

func TestResetChangelogSkipALevelShouldFail(t *testing.T) {
	s := NewChangelog()
	_, err := s.ResetTo(ReleaseHeading, "[Unreleased]")
	require.Error(t, err)
}

func TestResetChangelogToInvalidHeadingShouldFail(t *testing.T) {
	s := NewChangelog()
	s.ResetTo(IntroductionHeading, "title")
	_, err := s.ResetTo(ReleaseHeading, "Unreleased")
	require.Error(t, err)
}

func TestResetFilledChangelogToZeroIsSameAsPushToEmptyChangelog(t *testing.T) {
	require := require.New(t)

	s := NewChangelog()
	s.ResetTo(IntroductionHeading, "title")
	s.ResetTo(ReleaseHeading, "[Unreleased]")
	s.ResetTo(ChangeHeading, "Added")
	require.Equal("{title}{[Unreleased]}{Added}", s.String())

	_, err := s.ResetTo(IntroductionHeading, "other title")
	require.NoError(err)
	require.Equal("{other title}", s.String())
}

func TestResetFilledChangelogToOneIsSameAsTwoPushToEmptyChangelog(t *testing.T) {
	require := require.New(t)

	s := NewChangelog()
	s.ResetTo(IntroductionHeading, "title")
	s.ResetTo(ReleaseHeading, "[Unreleased]")
	s.ResetTo(ChangeHeading, "Added")

	require.Equal("{title}{[Unreleased]}{Added}", s.String())

	_, err := s.ResetTo(ReleaseHeading, "1.2.3 - 2020-04-15 [YANKED]")
	require.NoError(err)
	require.Equal("{title}{1.2.3 - 2020-04-15 [YANKED]}", s.String())
}
