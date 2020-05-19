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
