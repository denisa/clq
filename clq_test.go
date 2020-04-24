package main

import (
	"testing"
)

func TestChangelogShouldSucceed(t *testing.T) {
	executeClq(t, 0, "CHANGELOG.md")
}

func TestGoodDevWithDevModeShouldSucceed(t *testing.T) {
	executeClq(t, 0, "testdata/good_dev.md")
}

func TestGoodDevWithReleaseModeShouldFail(t *testing.T) {
	executeClq(t, 1, "-release", "testdata/good_dev.md")
}

func TestGoodRelWithDevModeShouldSucceed(t *testing.T) {
	executeClq(t, 0, "testdata/good_rel.md")
}

func TestGoodRelWithReleaseModeShouldSucceed(t *testing.T) {
	executeClq(t, 0, "-release", "testdata/good_rel.md")
}

func TestBadDevWithDevModeShouldFail(t *testing.T) {
	executeClq(t, 1, "testdata/bad_dev.md")
}

func TestBadDevWithReleaseModeShouldFail(t *testing.T) {
	executeClq(t, 1, "-release", "testdata/bad_dev.md")
}

func executeClq(t *testing.T, expected int, arguments ...string) {
	var actual = entryPoint("clq", arguments...)
	if expected != actual {
		t.Errorf("Expected %d but received %v", expected, actual)
	}

}
