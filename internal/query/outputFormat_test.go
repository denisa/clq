package query

import (
	"github.com/denisa/clq/internal/changelog"
)

func formatNoOutputDefined(format string) string {
	of, _ := newOutputFormat(format)
	of.open(newHeading(changelog.IntroductionHeading, "Changelog"))
	of.close(newHeading(changelog.IntroductionHeading, "Changelog"))
	return of.Result()
}

func formatIntroductionHeading(format string) string {
	of, _ := newOutputFormat(format)
	title := "Changelog"
	h := newHeading(changelog.IntroductionHeading, title)
	of.open(h)
	of.(resultCollector).setField("title", title)
	of.close(h)
	return of.Result()
}

func formatReleaseHeading(format string) string {
	of, _ := newOutputFormat(format)
	title := "[1.2.3] - 2020-05-16"
	h := newHeading(changelog.ReleaseHeading, title)
	of.open(h)
	of.(resultCollector).setField("title", title)
	of.close(h)
	return of.Result()
}

func formatChangeHeading(format string) string {
	of, _ := newOutputFormat(format)
	title := "Added"
	h := newHeading(changelog.ChangeHeading, title)
	of.open(h)
	of.(resultCollector).setField("title", title)
	of.close(h)
	return of.Result()
}

func formatChangeDescription(format string) string {
	of, _ := newOutputFormat(format)
	title := "foo"
	h := newHeading(changelog.ChangeDescription, title)
	of.open(h)
	of.(resultCollector).set(title)
	of.close(h)
	return of.Result()
}

func formatLoneArray(format string) string {
	of, _ := newOutputFormat(format)

	h := newHeading(changelog.ChangeHeading, "Added")
	of.open(h)
	of.(resultCollector).array("changes")
	{
		h := newHeading(changelog.ChangeDescription, "ignored")
		of.open(h)
		of.(resultCollector).set("foo")
		of.close(h)
	}
	{
		h := newHeading(changelog.ChangeDescription, "ignored")
		of.open(h)
		of.(resultCollector).set("bar")
		of.close(h)
	}
	of.close(h)
	return of.Result()
}
