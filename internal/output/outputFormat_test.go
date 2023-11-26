package output

import (
	"testing"

	"github.com/denisa/clq/internal/changelog"

	"github.com/stretchr/testify/require"
)

func TestUnsupportedOutputFormat(t *testing.T) {
	_, err := NewOutputFormat("yaml")
	require.Error(t, err)
}

func formatNoOutputDefined(format string) string {
	of, _ := NewOutputFormat(format)
	of.Open(newHeading(changelog.IntroductionHeading, "Changelog"))
	of.Close(newHeading(changelog.IntroductionHeading, "Changelog"))
	return of.Result()
}

func formatIntroductionHeading(format string) string {
	of, _ := NewOutputFormat(format)
	title := "Changelog"
	h := newHeading(changelog.IntroductionHeading, title)
	of.Open(h)
	of.SetField("title", title)
	of.Close(h)
	return of.Result()
}

func formatReleaseHeading(format string) string {
	of, _ := NewOutputFormat(format)
	title := "[1.2.3] - 2020-05-16"
	h := newHeading(changelog.ReleaseHeading, title)
	of.Open(h)
	of.SetField("title", title)
	of.Close(h)
	return of.Result()
}

func formatChangeHeading(format string) string {
	of, _ := NewOutputFormat(format)
	title := "Added"
	h := newHeading(changelog.ChangeHeading, title)
	of.Open(h)
	of.SetField("title", title)
	of.Close(h)
	return of.Result()
}

func formatChangeDescription(format string) string {
	of, _ := NewOutputFormat(format)
	title := "foo"
	h := newHeading(changelog.ChangeDescription, title)
	of.Open(h)
	of.Set(title)
	of.Close(h)
	return of.Result()
}

func formatLoneArray(format string) string {
	of, _ := NewOutputFormat(format)

	h := newHeading(changelog.ChangeHeading, "Added")
	of.Open(h)
	of.Array("changes")
	{
		h := newHeading(changelog.ChangeDescription, "ignored")
		of.Open(h)
		of.Set("foo")
		of.Close(h)
	}
	{
		h := newHeading(changelog.ChangeDescription, "ignored")
		of.Open(h)
		of.Set("bar")
		of.Close(h)
	}
	of.Close(h)
	return of.Result()
}

func newHeading(kind changelog.HeadingKind, text string) changelog.Heading {
	ck, _ := changelog.NewChangeKind("")
	hf := changelog.NewHeadingFactory(ck)

	if h, err := hf.NewHeading(kind, text); err != nil {
		panic(err)
	} else {
		return h
	}
}
