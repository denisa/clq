package changelog

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/blang/semver"
	"github.com/stretchr/testify/require"
)

func TestNewHeadingRelease(t *testing.T) {
	h, _ := NewHeading(ReleaseHeading, "[Unreleased]")
	requireHeadingInterface(t, "[Unreleased]", h)
}

func TestReleaseUnreleased(t *testing.T) {
	h, _ := newRelease("[Unreleased]")
	r, _ := h.(Release)

	require := require.New(t)
	require.Empty(r.Version())
	require.Empty(r.Date())
	require.Empty(r.Label())
	require.False(r.IsRelease())
	require.False(r.HasBeenYanked())
	require.False(r.HasBeenReleased())
}

func TestReleaseReleasedNoLabel(t *testing.T) {
	h, _ := newRelease("[1.2.3] - 2020-04-15")
	r, _ := h.(Release)

	require := require.New(t)
	require.Equal("1.2.3", r.Version())
	require.Equal("2020-04-15", r.Date())
	require.Empty(r.Label())
	require.True(r.IsRelease())
	require.False(r.HasBeenYanked())
	require.True(r.HasBeenReleased())

}

func TestReleaseReleasedLabel(t *testing.T) {
	h, _ := newRelease("[1.2.3] - 2020-04-15 Espelho")
	r, _ := h.(Release)

	require := require.New(t)
	require.Equal("1.2.3", r.Version())
	require.Equal("2020-04-15", r.Date())
	require.Equal("Espelho", r.Label())
	require.True(r.IsRelease())
	require.False(r.HasBeenYanked())
	require.True(r.HasBeenReleased())
}

func TestReleaseNonIsoDateShouldFail(t *testing.T) {
	_, err := newRelease("[1.2.3] - 2020.02.15")
	require.Error(t, err)
}

func TestReleaseReleasedDateNotInCalendar(t *testing.T) {
	_, err := newRelease("[1.2.3] - 2020-02-30")
	require.Error(t, err)
}

func TestReleaseReleasedAndYankedShouldFail(t *testing.T) {
	_, err := newRelease("[1.2.3] - 2020-02-15 [YANKED]")
	require.Error(t, err)
}

func TestReleaseReleasedVersionShouldFail(t *testing.T) {
	_, err := newRelease("[1.02] - 2020-02-15")
	require.Error(t, err)
}

func TestReleaseYanked(t *testing.T) {
	h, _ := newRelease("1.2.3 - 2020-04-15 [YANKED]")
	r, _ := h.(Release)

	require := require.New(t)
	require.Equal("1.2.3", r.Version())
	require.Equal("2020-04-15", r.Date())
	require.Empty(r.Label())
	require.True(r.IsRelease())
	require.True(r.HasBeenYanked())
	require.True(r.HasBeenReleased())
}

func TestReleaseYankedDateShouldFail(t *testing.T) {
	_, err := newRelease("1.2.3 - 2020-02-30 [YANKED]")
	require.Error(t, err)
}

func TestReleaseYankedVersionShouldFail(t *testing.T) {
	_, err := newRelease("1.02 - 2020-02-15 [YANKED]")
	require.Error(t, err)
}

func TestReleaseVersionEquality(t *testing.T) {
	h, _ := newRelease("[1.2.3] - 2020-04-15")
	r, _ := h.(Release)

	require := require.New(t)
	{
		v, _ := semver.Make("1.2.3")
		require.True(r.ReleaseIs(v))
	}
	{
		v, _ := semver.Make("1.2.4")
		require.False(r.ReleaseIs(v))
	}
}

func TestReleaseOrdering(t *testing.T) {
	h, _ := newRelease("[1.2.4] - 2020-04-16")
	r1, _ := h.(Release)

	h, _ = newRelease("[1.2.3] - 2020-04-15")
	r2, _ := h.(Release)

	require := require.New(t)
	require.NoError(r1.IsNewerThan(r2))
	require.Error(r2.IsNewerThan(r1))
}

func TestReleaseOrderingSameDay(t *testing.T) {
	h, _ := newRelease("[1.2.4] - 2020-04-16")
	r1, _ := h.(Release)

	h, _ = newRelease("[1.2.3] - 2020-04-16")
	r2, _ := h.(Release)

	require := require.New(t)
	require.NoError(r1.IsNewerThan(r2))
	require.Error(r2.IsNewerThan(r1))
}

func TestReleaseOrderingSameVersionShouldFail(t *testing.T) {
	h, _ := newRelease("[1.2.4] - 2020-04-16")
	r1, _ := h.(Release)

	h, _ = newRelease("[1.2.4] - 2020-04-15")
	r2, _ := h.(Release)

	require := require.New(t)
	require.Error(r1.IsNewerThan(r2))
	require.Error(r2.IsNewerThan(r1))
}

func TestReleaseOrderingMixedUp(t *testing.T) {
	h, _ := newRelease("[1.2.4] - 2020-04-15")
	r1, _ := h.(Release)
	//
	h, _ = newRelease("[1.2.3] - 2020-04-16")
	r2, _ := h.(Release)

	require := require.New(t)
	require.Error(r2.IsNewerThan(r1))
	require.Error(r1.IsNewerThan(r2))
}

func TestNextRelease(t *testing.T) {
	testcases := []struct {
		changeMap ChangeMap
		expected  semver.Version
	}{
		{ChangeMap{"Added": true}, semver.Version{Major: 1, Minor: 0, Patch: 0}},
		{ChangeMap{"Changed": true}, semver.Version{Major: 0, Minor: 1, Patch: 0}},
		{ChangeMap{"Fixed": true}, semver.Version{Major: 0, Minor: 0, Patch: 1}},
		{ChangeMap{}, semver.Version{}},
	}
	for _, testcase := range testcases {
		t.Run(testcase.expected.String(), func(t *testing.T) {
			r := &Release{}
			require.Equal(t, testcase.expected, r.NextRelease(testcase.changeMap))
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
			require := require.New(t)
			re, err := regexp.Compile(testcase.exp)
			if err != nil {
				require.Nil(err)
				return
			}

			matches := re.FindStringSubmatch(testcase.input)

			value := subexp(re.SubexpNames(), matches, testcase.subexp)
			require.Equal(testcase.value, value)
		})
	}
}
