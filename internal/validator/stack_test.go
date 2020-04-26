package validator

import (
	"testing"
)

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
		t.Error("change mot expected")
	}
}

func TestOneLevelStackIsDocument(t *testing.T) {
	s := NewStack()
	s.push("first")
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
	s.push("first")
	s.push("second")
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
	s.push("first")
	s.push("second")
	s.push("third")
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
	expected := "an item"
	s.push(expected)
	actual, _ := s.pop()
	assertStringEquals(t, expected, actual)
}

func TestPopEmptyStackFails(t *testing.T) {
	s := NewStack()
	_, err := s.pop()
	assertErrorExpected(t, err)
}

func TestStackDepthGrowsAndShrink(t *testing.T) {
	s := NewStack()
	assertIntEquals(t, 0, s.depth())
	assertStringEquals(t, "", s.asPath())

	s.push("first")
	assertIntEquals(t, 1, s.depth())
	assertStringEquals(t, "{first}", s.asPath())

	s.push("second")
	assertIntEquals(t, 2, s.depth())
	assertStringEquals(t, "{first}{second}", s.asPath())

	if actual, err := s.pop(); err == nil {
		assertStringEquals(t, "second", actual)
	} else {
		t.Error(err)
	}
	assertIntEquals(t, 1, s.depth())
	assertStringEquals(t, "{first}", s.asPath())

	if actual, err := s.pop(); err == nil {
		assertStringEquals(t, "first", actual)
	} else {
		t.Error(err)
	}

	assertIntEquals(t, 0, s.depth())
	assertStringEquals(t, "", s.asPath())
}

func TestResetEmptyStackToZeroIsSameAsPush(t *testing.T) {
	s := NewStack()
	s.resetTo(0, "first")
	assertStringEquals(t, "{first}", s.asPath())
}

func TestResetEmptyStackToALevelShouldFail(t *testing.T) {
	s := NewStack()
	err := s.resetTo(2, "first")
	assertErrorExpected(t, err)
}

func TestResetFilledStackToZeroIsSameAsPushToEmptyStack(t *testing.T) {
	s := NewStack()
	s.push("one")
	s.push("two")
	s.push("three")
	assertStringEquals(t, "{one}{two}{three}", s.asPath())
	s.resetTo(0, "first")
	assertStringEquals(t, "{first}", s.asPath())
}

func TestResetFilledStackToOneIsSameAsTwoPushToEmptyStack(t *testing.T) {
	s := NewStack()
	s.push("one")
	s.push("two")
	s.push("three")
	assertStringEquals(t, "{one}{two}{three}", s.asPath())
	s.resetTo(1, "first")
	assertStringEquals(t, "{one}{first}", s.asPath())
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
