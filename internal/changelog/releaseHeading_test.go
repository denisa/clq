package changelog

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/blang/semver/v4"
	semverConstants "github.com/denisa/clq/internal/semver"
	"github.com/stretchr/testify/require"
)

func TestNewHeadingRelease(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	h, _ := hf.NewHeading(ReleaseHeading, "[Unreleased]")
	requireHeadingInterface(t, "[Unreleased]", h)
}

func TestReleaseUnreleased(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	h, _ := hf.newRelease("[Unreleased]")
	r, _ := h.(Release)

	assertions := require.New(t)
	assertions.Empty(r.Version())
	assertions.Empty(r.Date())
	assertions.Empty(r.Label())
	assertions.False(r.IsRelease())
	assertions.False(r.IsPrerelease())
	assertions.False(r.HasBeenYanked())
	assertions.False(r.HasBeenReleased())
}

func TestReleasePrereleasedNoLabel(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	h, _ := hf.newRelease("[1.2.3-rc.1] - 2020-04-15")
	r, _ := h.(Release)

	assertions := require.New(t)
	assertions.Equal("1.2.3-rc.1", r.Version())
	assertions.Equal("2020-04-15", r.Date())
	assertions.Empty(r.Label())
	assertions.False(r.IsRelease())
	assertions.True(r.IsPrerelease())
	assertions.False(r.HasBeenYanked())
	assertions.True(r.HasBeenReleased())
}

func TestReleaseReleasedNoLabel(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	h, _ := hf.newRelease("[1.2.3] - 2020-04-15")
	r, _ := h.(Release)

	assertions := require.New(t)
	assertions.Equal("1.2.3", r.Version())
	assertions.Equal("2020-04-15", r.Date())
	assertions.Empty(r.Label())
	assertions.True(r.IsRelease())
	assertions.False(r.IsPrerelease())
	assertions.False(r.HasBeenYanked())
	assertions.True(r.HasBeenReleased())
}

func TestReleaseReleasedLabel(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	h, _ := hf.newRelease("[1.2.3] - 2020-04-15 Espelho")
	r, _ := h.(Release)

	assertions := require.New(t)
	assertions.Equal("1.2.3", r.Version())
	assertions.Equal("2020-04-15", r.Date())
	assertions.Equal("Espelho", r.Label())
	assertions.True(r.IsRelease())
	assertions.False(r.IsPrerelease())
	assertions.False(r.HasBeenYanked())
	assertions.True(r.HasBeenReleased())
}

func TestReleaseNonIsoDateSeparatorShouldFail(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	_, err := hf.newRelease("[1.2.3] - 2020.02.15")
	require.Error(t, err)
}

func TestReleaseNonIsoDateSingleDigitShouldFail(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	_, err := hf.newRelease("[1.2.3] - 2020-4-1")
	require.Error(t, err)
}

func TestReleaseReleasedDateNotInCalendar(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	_, err := hf.newRelease("[1.2.3] - 2020-02-30")
	require.Error(t, err)
}

func TestReleaseReleasedAndYankedShouldFail(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	_, err := hf.newRelease("[1.2.3] - 2020-02-15 [YANKED]")
	require.Error(t, err)
}

func TestReleaseReleasedVersionShouldFail(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	_, err := hf.newRelease("[1.02] - 2020-02-15")
	require.Error(t, err)
}

func TestReleaseYanked(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	h, _ := hf.newRelease("1.2.3 - 2020-04-15 [YANKED]")
	r, _ := h.(Release)

	assertions := require.New(t)
	assertions.Equal("1.2.3", r.Version())
	assertions.Equal("2020-04-15", r.Date())
	assertions.Empty(r.Label())
	assertions.True(r.IsRelease())
	assertions.True(r.HasBeenYanked())
	assertions.True(r.HasBeenReleased())
}

func TestReleaseYankedDateShouldFail(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	_, err := hf.newRelease("1.2.3 - 2020-02-30 [YANKED]")
	require.Error(t, err)
}

func TestReleaseYankedVersionShouldFail(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	_, err := hf.newRelease("1.02 - 2020-02-15 [YANKED]")
	require.Error(t, err)
}

func TestReleaseVersionEquality(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	h, _ := hf.newRelease("[1.2.3] - 2020-04-15")
	r, _ := h.(Release)

	assertions := require.New(t)
	{
		v, _ := semver.Make("1.2.3")
		assertions.True(r.ReleaseIs(v))
	}
	{
		v, _ := semver.Make("1.2.4")
		assertions.False(r.ReleaseIs(v))
	}
}

func TestReleaseOrdering(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	h, _ := hf.newRelease("[1.2.4] - 2020-04-16")
	r1, _ := h.(Release)

	h, _ = hf.newRelease("[1.2.3] - 2020-04-15")
	r2, _ := h.(Release)

	assertions := require.New(t)
	assertions.NoError(r1.IsNewerThan(r2))
	assertions.Error(r2.IsNewerThan(r1))
}

func TestReleaseOrderingSameDay(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	h, _ := hf.newRelease("[1.2.4] - 2020-04-16")
	r1, _ := h.(Release)

	h, _ = hf.newRelease("[1.2.3] - 2020-04-16")
	r2, _ := h.(Release)

	assertions := require.New(t)
	assertions.NoError(r1.IsNewerThan(r2))
	assertions.Error(r2.IsNewerThan(r1))
}

func TestReleaseOrderingSameVersionShouldFail(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	h, _ := hf.newRelease("[1.2.4] - 2020-04-16")
	r1, _ := h.(Release)

	h, _ = hf.newRelease("[1.2.4] - 2020-04-15")
	r2, _ := h.(Release)

	assertions := require.New(t)
	assertions.Error(r1.IsNewerThan(r2))
	assertions.Error(r2.IsNewerThan(r1))
}

func TestReleaseOrderingMixedUp(t *testing.T) {
	ck, _ := NewChangeKind("")
	hf := NewHeadingFactory(ck)
	h, _ := hf.newRelease("[1.2.4] - 2020-04-15")
	r1, _ := h.(Release)
	//
	h, _ = hf.newRelease("[1.2.3] - 2020-04-16")
	r2, _ := h.(Release)

	assertions := require.New(t)
	assertions.Error(r2.IsNewerThan(r1))
	assertions.Error(r1.IsNewerThan(r2))
}

func TestNextRelease(t *testing.T) {
	testcases := []struct {
		semverIdentifier semverConstants.Identifier
		expected         semver.Version
	}{
		{semverConstants.Major, semver.Version{Major: 1, Minor: 0, Patch: 0}},
		{semverConstants.Minor, semver.Version{Major: 0, Minor: 1, Patch: 0}},
		{semverConstants.Patch, semver.Version{Major: 0, Minor: 0, Patch: 1}},
		{semverConstants.Prerelease, semver.Version{}},
		{semverConstants.Build, semver.Version{}},
	}
	for _, testcase := range testcases {
		t.Run(testcase.expected.String(), func(t *testing.T) {
			r := &Release{}
			require.Equal(t, testcase.expected, r.NextRelease(testcase.semverIdentifier))
		})
	}
}

func TestSubexp(t *testing.T) {
	testcases := []struct {
		input  string
		exp    string
		subexp string
		value  string
	}{
		{
			`a b c`,
			`(?P<foo>.)\s+(?P<blah>.*)`,
			`foo`,
			`a`,
		},
		{
			`a b c`,
			`(?P<foo>.)\s+(?P<blah>.*)`,
			`blah`,
			`b c`,
		},
		{
			`a b c`,
			`(?P<foo>.)\s+(?P<blah>.*)`,
			`unknown`,
			``,
		},
		{
			`a b c`,
			`(?P<foo>1)?a`,
			`foo`,
			``,
		},
	}

	for index, testcase := range testcases {
		t.Run(strconv.Itoa(index), func(t *testing.T) {
			assertions := require.New(t)
			re, err := regexp.Compile(testcase.exp)
			if err != nil {
				assertions.Nil(err)
				return
			}

			matches := re.FindStringSubmatch(testcase.input)

			value := subexp(re.SubexpNames(), matches, testcase.subexp)
			assertions.Equal(testcase.value, value)
		})
	}
}
