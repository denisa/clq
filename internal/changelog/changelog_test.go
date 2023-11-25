package changelog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewChangelogIsEmpty(t *testing.T) {
	require := require.New(t)

	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	s := NewChangelog(hf)

	require.False(s.Introduction(), "title not expected")
	require.False(s.Release(), "release not expected")
	require.False(s.Change(), "change not expected")
}

func TestOneLevelChangelogIsDocument(t *testing.T) {
	require := require.New(t)

	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	s := NewChangelog(hf)
	s.Section(IntroductionHeading, "title")

	require.True(s.Introduction(), "title expected")
	require.False(s.Release(), "release not expected")
	require.False(s.Change(), "change not expected")
}

func TestTwoLevelChangelogIsRelease(t *testing.T) {
	require := require.New(t)

	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	s := NewChangelog(hf)
	s.Section(IntroductionHeading, "title")
	s.Section(ReleaseHeading, "[Unreleased]")

	require.False(s.Introduction(), "title not expected")
	require.True(s.Release(), "release expected")
	require.False(s.Change(), "change not expected")

}

func TestThreeLevelChangelogIsChange(t *testing.T) {
	require := require.New(t)

	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	s := NewChangelog(hf)
	s.Section(IntroductionHeading, "title")
	s.Section(ReleaseHeading, "[Unreleased]")
	s.Section(ChangeHeading, "Added")

	require.False(s.Introduction(), "title not expected")
	require.False(s.Release(), "release not expected")
	require.True(s.Change(), "change expected")
}

func TestResetChangelogSkipALevelShouldFail(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	s := NewChangelog(hf)
	_, err := s.Section(ReleaseHeading, "[Unreleased]")
	require.Error(t, err)
}

func TestResetChangelogToInvalidHeadingShouldFail(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	s := NewChangelog(hf)
	s.Section(IntroductionHeading, "title")
	_, err := s.Section(ReleaseHeading, "Unreleased")
	require.Error(t, err)
}

func TestResetFilledChangelogToZeroIsSameAsPushToEmptyChangelog(t *testing.T) {
	require := require.New(t)

	recorder := &recorder{}

	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	s := NewChangelog(hf)
	s.Listener(recorder)

	s.Section(IntroductionHeading, "title")
	requireEventsEquals(require, &[]string{"Enter {title}"}, &recorder.events)

	s.Section(ReleaseHeading, "[Unreleased]")
	requireEventsEquals(require, &[]string{"Enter {title}", "Enter {[Unreleased]}"}, &recorder.events)

	s.Section(ChangeHeading, "Added")
	requireEventsEquals(require, &[]string{"Enter {title}", "Enter {[Unreleased]}", "Enter {Added}"}, &recorder.events)
	require.Equal("{title}{[Unreleased]}{Added}", s.String())

	_, err := s.Section(IntroductionHeading, "other title")
	require.NoError(err)
	requireEventsEquals(require, &[]string{"Enter {title}", "Enter {[Unreleased]}", "Enter {Added}", "Exit {Added}", "Exit {[Unreleased]}", "Exit {title}", "Enter {other title}"}, &recorder.events)
	require.Equal("{other title}", s.String())
}

func TestResetFilledChangelogToOneIsSameAsTwoPushToEmptyChangelog(t *testing.T) {
	require := require.New(t)

	recorder := &recorder{}

	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	s := NewChangelog(hf)
	s.Listener(recorder)

	s.Section(IntroductionHeading, "title")
	requireEventsEquals(require, &[]string{"Enter {title}"}, &recorder.events)

	s.Section(ReleaseHeading, "[Unreleased]")
	requireEventsEquals(require, &[]string{"Enter {title}", "Enter {[Unreleased]}"}, &recorder.events)

	s.Section(ChangeHeading, "Added")
	requireEventsEquals(require, &[]string{"Enter {title}", "Enter {[Unreleased]}", "Enter {Added}"}, &recorder.events)
	require.Equal("{title}{[Unreleased]}{Added}", s.String())

	_, err := s.Section(ReleaseHeading, "1.2.3 - 2020-04-15 [YANKED]")
	require.NoError(err)
	requireEventsEquals(require, &[]string{"Enter {title}", "Enter {[Unreleased]}", "Enter {Added}", "Exit {Added}", "Exit {[Unreleased]}", "Enter {1.2.3 - 2020-04-15 [YANKED]}"}, &recorder.events)
	require.Equal("{title}{1.2.3 - 2020-04-15 [YANKED]}", s.String())
}

func requireEventsEquals(require *require.Assertions, expected *[]string, actual *[]string) {
	require.Equalf(len(*expected), len(*actual), "%v\n%v", *expected, *actual)
	for i := range *expected {
		require.Equal((*expected)[i], (*actual)[i])
	}
}
