package validator

import (
	"reflect"
	"regexp"
	"strconv"
	"testing"

	"github.com/blang/semver"
)

func TestTitle(t *testing.T) {
	h, _ := newTitle("changelog")
	assertHeadingInterface(t, "changelog", h)
}

func TestRelease(t *testing.T) {
	h, _ := newRelease("[Unreleased]")
	assertHeadingInterface(t, "[Unreleased]", h)
}

func TestChange(t *testing.T) {
	h, _ := newChange("Security")
	assertHeadingInterface(t, "Security", h)
}

func TestNextRelease(t *testing.T) {
	assertVersionEquals(t, semver.Version{Major: 1, Minor: 0, Patch: 0}, nextRelease(changeMap{"Added": true}, semver.Version{}))
	assertVersionEquals(t, semver.Version{Major: 0, Minor: 1, Patch: 0}, nextRelease(changeMap{"Changed": true}, semver.Version{}))
	assertVersionEquals(t, semver.Version{Major: 0, Minor: 0, Patch: 1}, nextRelease(changeMap{"Fixed": true}, semver.Version{}))
	assertVersionEquals(t, semver.Version{Major: 0, Minor: 0, Patch: 0}, nextRelease(changeMap{}, semver.Version{}))
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
			re, err := regexp.Compile(testcase.exp)
			if err != nil {
				assertNoErrorExpected(t, err)
				return
			}

			matches := re.FindStringSubmatch(testcase.input)

			value := subexp(re.SubexpNames(), matches, testcase.subexp)
			assertStringEquals(t, testcase.value, value)
		})
	}
}

func assertHeadingInterface(t *testing.T, name string, actual Heading) {
	assertStringEquals(t, name, actual.Name())
	assertStringEquals(t, asPath(actual.Name()), actual.AsPath())
}

func assertHeadingEquals(t *testing.T, expected Heading, actual Heading) {
	if expected.Name() != actual.Name() {
		t.Errorf("Expected %v but received %v", expected, actual)
	}
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		t.Errorf("Expected %v but received %v", reflect.TypeOf(expected), reflect.TypeOf(actual))
	}

}

func assertVersionEquals(t *testing.T, expected semver.Version, actual semver.Version) {
	assertStringEquals(t, expected.String(), actual.String())
}
