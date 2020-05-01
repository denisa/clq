package validator

import (
	"reflect"
	"regexp"
	"strconv"
	"testing"
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
	s := NewStack()
	if !s.empty() {
		t.Error("empty stack expected")
	}
	if s.title() {
		t.Error("title not expected")
	}
	if s.release() {
		t.Error("release not expected")
	}
	if s.change() {
		t.Error("change not expected")
	}
}

func TestOneLevelStackIsDocument(t *testing.T) {
	s := NewStack()
	s.push(newTestHeading("first"))
	if s.empty() {
		t.Error("empty stack not expected")
	}
	if !s.title() {
		t.Error("title expected")
	}
	if s.release() {
		t.Error("release not expected")
	}
	if s.change() {
		t.Error("change not expected")
	}
}
func TestTwoLevelStackIsRelease(t *testing.T) {
	s := NewStack()
	s.push(newTestHeading("first"))
	s.push(newTestHeading("second"))
	if s.empty() {
		t.Error("empty stack not expected")
	}
	if s.title() {
		t.Error("title not expected")
	}
	if !s.release() {
		t.Error("release expected")
	}
	if s.change() {
		t.Error("change not expected")
	}
}
func TestThreeLevelStackIsChange(t *testing.T) {
	s := NewStack()
	s.push(newTestHeading("first"))
	s.push(newTestHeading("second"))
	s.push(newTestHeading("third"))
	if s.empty() {
		t.Error("empty stack not expected")
	}
	if s.title() {
		t.Error("title expected")
	}
	if s.release() {
		t.Error("release not expected")
	}
	if !s.change() {
		t.Error("change expected")
	}
}

func TestPopReturnsPushedContent(t *testing.T) {
	s := NewStack()
	expected := newTestHeading("an item")
	s.push(expected)
	assertIntEquals(t, 1, s.depth())

	actual, _ := s.pop()
	assertHeadingEquals(t, expected, actual)
	assertIntEquals(t, 0, s.depth())
}

func TestPopEmptyStackFails(t *testing.T) {
	s := NewStack()
	_, err := s.pop()
	assertErrorExpected(t, err)
}
func TestPeekReturnsPushedContent(t *testing.T) {
	s := NewStack()
	expected := newTestHeading("an item")
	s.push(expected)
	assertIntEquals(t, 1, s.depth())

	actual, _ := s.Peek()
	assertHeadingEquals(t, expected, actual)
	assertIntEquals(t, 1, s.depth())
}

func TestPeekEmptyStackFails(t *testing.T) {
	s := NewStack()
	_, err := s.Peek()
	assertErrorExpected(t, err)
}

func TestStackDepthGrowsAndShrink(t *testing.T) {
	s := NewStack()
	assertIntEquals(t, 0, s.depth())
	assertStringEquals(t, "", s.AsPath())

	s.push(newTestHeading("first"))
	assertIntEquals(t, 1, s.depth())
	assertStringEquals(t, "{first}", s.AsPath())

	s.push(newTestHeading("second"))
	assertIntEquals(t, 2, s.depth())
	assertStringEquals(t, "{first}{second}", s.AsPath())

	if actual, err := s.pop(); err == nil {
		assertHeadingEquals(t, newTestHeading("second"), actual)
	} else {
		t.Error(err)
	}
	assertIntEquals(t, 1, s.depth())
	assertStringEquals(t, "{first}", s.AsPath())

	if actual, err := s.pop(); err == nil {
		assertHeadingEquals(t, newTestHeading("first"), actual)
	} else {
		t.Error(err)
	}

	assertIntEquals(t, 0, s.depth())
	assertStringEquals(t, "", s.AsPath())
}

func TestResetEmptyStackToZeroIsSameAsPush(t *testing.T) {
	s := NewStack()
	s.ResetTo(titleHeading, "title")
	assertStringEquals(t, "{title}", s.AsPath())
}

func TestResetEmptyStackToALevelShouldFail(t *testing.T) {
	s := NewStack()
	err := s.ResetTo(releaseHeading, "[Unreleased]")
	assertErrorExpected(t, err)
}

func TestResetFilledStackToZeroIsSameAsPushToEmptyStack(t *testing.T) {
	s := NewStack()
	if h, err := newTitle("title"); err == nil {
		s.push(h)
	}
	if h, err := newRelease("[Unreleased]"); err == nil {
		s.push(h)
	}
	if h, err := newChange("Added"); err == nil {
		s.push(h)
	}
	assertStringEquals(t, "{title}{[Unreleased]}{Added}", s.AsPath())

	s.ResetTo(titleHeading, "other title")
	assertStringEquals(t, "{other title}", s.AsPath())
}

func TestResetFilledStackToOneIsSameAsTwoPushToEmptyStack(t *testing.T) {
	s := NewStack()
	if h, err := newTitle("title"); err == nil {
		s.push(h)
	}
	if h, err := newRelease("[Unreleased]"); err == nil {
		s.push(h)
	}
	if h, err := newChange("Added"); err == nil {
		s.push(h)
	}
	assertStringEquals(t, "{title}{[Unreleased]}{Added}", s.AsPath())

	s.ResetTo(releaseHeading, "1.2.3 - 2020-04-15 [YANKED]")
	assertStringEquals(t, "{title}{1.2.3 - 2020-04-15 [YANKED]}", s.AsPath())
}

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

			value := subexp(re, matches, testcase.subexp)
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

func assertIntEquals(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Errorf("Expected %d but received %d", expected, actual)
	}
}

func assertStringEquals(t *testing.T, expected string, actual string) {
	if expected != actual {
		t.Errorf("Expected %v but received %v", expected, actual)
	}
}

func assertNoErrorExpected(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Expected no error, received %v", err)
	}
}

func assertErrorExpected(t *testing.T, err error) {
	if err == nil {
		t.Error("Expected an error")
	}
}
