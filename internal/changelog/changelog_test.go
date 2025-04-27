package changelog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewChangelogIsEmpty(t *testing.T) {
	assertions := require.New(t)

	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	s := NewChangelog(hf)

	assertions.False(s.Introduction(), "title not expected")
	assertions.False(s.Release(), "release not expected")
	assertions.False(s.Change(), "change not expected")
}

func TestOneLevelChangelogIsDocument(t *testing.T) {
	assertions := require.New(t)

	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	s := NewChangelog(hf)
	_, _ = s.Section(IntroductionHeading, "title")

	assertions.True(s.Introduction(), "title expected")
	assertions.False(s.Release(), "release not expected")
	assertions.False(s.Change(), "change not expected")
}

func TestTwoLevelChangelogIsRelease(t *testing.T) {
	assertions := require.New(t)

	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	s := NewChangelog(hf)
	_, _ = s.Section(IntroductionHeading, "title")
	_, _ = s.Section(ReleaseHeading, "[Unreleased]")

	assertions.False(s.Introduction(), "title not expected")
	assertions.True(s.Release(), "release expected")
	assertions.False(s.Change(), "change not expected")

}

func TestThreeLevelChangelogIsChange(t *testing.T) {
	assertions := require.New(t)

	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	s := NewChangelog(hf)
	_, _ = s.Section(IntroductionHeading, "title")
	_, _ = s.Section(ReleaseHeading, "[Unreleased]")
	_, _ = s.Section(ChangeHeading, "Added")

	assertions.False(s.Introduction(), "title not expected")
	assertions.False(s.Release(), "release not expected")
	assertions.True(s.Change(), "change expected")
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
	_, _ = s.Section(IntroductionHeading, "title")
	_, err := s.Section(ReleaseHeading, "Unreleased")
	require.Error(t, err)
}

func TestResetFilledChangelogToZeroIsSameAsPushToEmptyChangelog(t *testing.T) {
	assertions := require.New(t)

	recorder := &recorder{}

	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	s := NewChangelog(hf)
	s.Listener(recorder)

	_, _ = s.Section(IntroductionHeading, "title")
	requireEventsEquals(assertions, &[]string{"Enter {title}"}, &recorder.events)

	_, _ = s.Section(ReleaseHeading, "[Unreleased]")
	requireEventsEquals(assertions, &[]string{"Enter {title}", "Enter {[Unreleased]}"}, &recorder.events)

	_, _ = s.Section(ChangeHeading, "Added")
	requireEventsEquals(assertions, &[]string{"Enter {title}", "Enter {[Unreleased]}", "Enter {Added}"}, &recorder.events)
	assertions.Equal("{title}{[Unreleased]}{Added}", s.String())

	_, err := s.Section(IntroductionHeading, "other title")
	assertions.NoError(err)
	requireEventsEquals(assertions, &[]string{"Enter {title}", "Enter {[Unreleased]}", "Enter {Added}", "Exit {Added}", "Exit {[Unreleased]}", "Exit {title}", "Enter {other title}"}, &recorder.events)
	assertions.Equal("{other title}", s.String())
}

func TestResetFilledChangelogToOneIsSameAsTwoPushToEmptyChangelog(t *testing.T) {
	assertions := require.New(t)

	recorder := &recorder{}

	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	s := NewChangelog(hf)
	s.Listener(recorder)

	_, _ = s.Section(IntroductionHeading, "title")
	requireEventsEquals(assertions, &[]string{"Enter {title}"}, &recorder.events)

	_, _ = s.Section(ReleaseHeading, "[Unreleased]")
	requireEventsEquals(assertions, &[]string{"Enter {title}", "Enter {[Unreleased]}"}, &recorder.events)

	_, _ = s.Section(ChangeHeading, "Added")
	requireEventsEquals(assertions, &[]string{"Enter {title}", "Enter {[Unreleased]}", "Enter {Added}"}, &recorder.events)
	assertions.Equal("{title}{[Unreleased]}{Added}", s.String())

	_, err := s.Section(ReleaseHeading, "[1.2.3] - 2020-04-15 [YANKED]")
	assertions.NoError(err)
	requireEventsEquals(assertions, &[]string{"Enter {title}", "Enter {[Unreleased]}", "Enter {Added}", "Exit {Added}", "Exit {[Unreleased]}", "Enter {[1.2.3] - 2020-04-15 [YANKED]}"}, &recorder.events)
	assertions.Equal("{title}{[1.2.3] - 2020-04-15 [YANKED]}", s.String())
}

func requireEventsEquals(assertions *require.Assertions, expected *[]string, actual *[]string) {
	assertions.Equalf(len(*expected), len(*actual), "%v\n%v", *expected, *actual)
	for i := range *expected {
		assertions.Equal((*expected)[i], (*actual)[i])
	}
}
