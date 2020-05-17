package changelog

import (
	"reflect"
	"regexp"
	"strconv"
	"testing"

	"github.com/blang/semver"
	"github.com/stretchr/testify/require"
)

func TestNewHeadingChangelog(t *testing.T) {
	h, _ := NewHeading(TitleHeading, "changelog")
	requireHeadingInterface(t, "changelog", h)
}

func TestNewHeadingRelease(t *testing.T) {
	h, _ := NewHeading(ReleaseHeading, "[Unreleased]")
	requireHeadingInterface(t, "[Unreleased]", h)
}

func TestNewHeadingChange(t *testing.T) {
	h, _ := NewHeading(ChangeHeading, "Security")
	requireHeadingInterface(t, "Security", h)
}

func TestNewHeadingUnknown(t *testing.T) {
	_, err := NewHeading(-1, "Who knows what")
	require.Error(t, err)
}

func TestChangelog(t *testing.T) {
	h, _ := newChangelog("changelog")
	requireHeadingInterface(t, "changelog", h)
}

func TestRelease(t *testing.T) {
	h, _ := newRelease("[Unreleased]")
	requireHeadingInterface(t, "[Unreleased]", h)
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

func TestChange(t *testing.T) {
	h, _ := newChange("Security")
	requireHeadingInterface(t, "Security", h)
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

func requireHeadingInterface(t *testing.T, name string, actual Heading) {
	require.Equal(t, name, actual.Name())
	require.Equal(t, asPath(actual.Name()), actual.AsPath())
}

func requireHeadingEquals(require *require.Assertions, expected Heading, actual Heading) {
	require.Equal(expected, actual)
	require.Equal(reflect.TypeOf(expected), reflect.TypeOf(actual))
}
