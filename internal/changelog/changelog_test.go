package changelog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewChangelogIsEmpty(t *testing.T) {
	require := require.New(t)

	s := NewChangelog()

	require.False(s.Introduction(), "title not expected")
	require.False(s.Release(), "release not expected")
	require.False(s.Change(), "change not expected")
}

func TestOneLevelChangelogIsDocument(t *testing.T) {
	require := require.New(t)

	s := NewChangelog()
	s.Section(IntroductionHeading, "title")

	require.True(s.Introduction(), "title expected")
	require.False(s.Release(), "release not expected")
	require.False(s.Change(), "change not expected")
}

func TestTwoLevelChangelogIsRelease(t *testing.T) {
	require := require.New(t)

	s := NewChangelog()
	s.Section(IntroductionHeading, "title")
	s.Section(ReleaseHeading, "[Unreleased]")

	require.False(s.Introduction(), "title not expected")
	require.True(s.Release(), "release expected")
	require.False(s.Change(), "change not expected")

}

func TestThreeLevelChangelogIsChange(t *testing.T) {
	require := require.New(t)

	s := NewChangelog()
	s.Section(IntroductionHeading, "title")
	s.Section(ReleaseHeading, "[Unreleased]")
	s.Section(ChangeHeading, "Added")

	require.False(s.Introduction(), "title not expected")
	require.False(s.Release(), "release not expected")
	require.True(s.Change(), "change expected")
}

func TestResetChangelogSkipALevelShouldFail(t *testing.T) {
	s := NewChangelog()
	_, err := s.Section(ReleaseHeading, "[Unreleased]")
	require.Error(t, err)
}

func TestResetChangelogToInvalidHeadingShouldFail(t *testing.T) {
	s := NewChangelog()
	s.Section(IntroductionHeading, "title")
	_, err := s.Section(ReleaseHeading, "Unreleased")
	require.Error(t, err)
}

func TestResetFilledChangelogToZeroIsSameAsPushToEmptyChangelog(t *testing.T) {
	require := require.New(t)

	s := NewChangelog()
	s.Section(IntroductionHeading, "title")
	s.Section(ReleaseHeading, "[Unreleased]")
	s.Section(ChangeHeading, "Added")
	require.Equal("{title}{[Unreleased]}{Added}", s.String())

	_, err := s.Section(IntroductionHeading, "other title")
	require.NoError(err)
	require.Equal("{other title}", s.String())
}

func TestResetFilledChangelogToOneIsSameAsTwoPushToEmptyChangelog(t *testing.T) {
	require := require.New(t)

	s := NewChangelog()
	s.Section(IntroductionHeading, "title")
	s.Section(ReleaseHeading, "[Unreleased]")
	s.Section(ChangeHeading, "Added")

	require.Equal("{title}{[Unreleased]}{Added}", s.String())

	_, err := s.Section(ReleaseHeading, "1.2.3 - 2020-04-15 [YANKED]")
	require.NoError(err)
	require.Equal("{title}{1.2.3 - 2020-04-15 [YANKED]}", s.String())
}
