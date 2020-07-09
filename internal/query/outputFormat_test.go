package query

import (
	"github.com/denisa/clq/internal/changelog"
)

func formatNoOutputDefined(format string) string {
	of, _ := newOutputFormat(format)
	of.Open(newHeading(changelog.IntroductionHeading, "Changelog"))
	of.Close(newHeading(changelog.IntroductionHeading, "Changelog"))
	return of.Result()
}

func formatIntroductionHeading(format string) string {
	of, _ := newOutputFormat(format)
	title := "Changelog"
	h := newHeading(changelog.IntroductionHeading, title)
	of.Open(h)
	of.(ResultCollector).setField("title", title)
	of.Close(h)
	return of.Result()
}

func formatReleaseHeading(format string) string {
	of, _ := newOutputFormat(format)
	title := "[1.2.3] - 2020-05-16"
	h := newHeading(changelog.ReleaseHeading, title)
	of.Open(h)
	of.(ResultCollector).setField("title", title)
	of.Close(h)
	return of.Result()
}

func formatChangeHeading(format string) string {
	of, _ := newOutputFormat(format)
	title := "Added"
	h := newHeading(changelog.ChangeHeading, title)
	of.Open(h)
	of.(ResultCollector).setField("title", title)
	of.Close(h)
	return of.Result()
}

func formatChangeDescription(format string) string {
	of, _ := newOutputFormat(format)
	title := "foo"
	h := newHeading(changelog.ChangeDescription, title)
	of.Open(h)
	of.(ResultCollector).set(title)
	of.Close(h)
	return of.Result()
}

func formatLoneArray(format string) string {
	of, _ := newOutputFormat(format)

	h := newHeading(changelog.ChangeHeading, "Added")
	of.Open(h)
	of.(ResultCollector).array("changes")
	{
		h := newHeading(changelog.ChangeDescription, "ignored")
		of.Open(h)
		of.(ResultCollector).set("foo")
		of.Close(h)
	}
	{
		h := newHeading(changelog.ChangeDescription, "ignored")
		of.Open(h)
		of.(ResultCollector).set("bar")
		of.Close(h)
	}
	of.Close(h)
	return of.Result()
}
