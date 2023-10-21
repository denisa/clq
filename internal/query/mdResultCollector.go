package query

import (
	"strings"

	"github.com/denisa/clq/internal/changelog"
)

// a mdResultCollector produces a markdown-repesentation of the query result.
type mdResultCollector struct {
	result strings.Builder
	prefix string
}

func (rc *mdResultCollector) Result() string {
	return strings.TrimSuffix(rc.result.String(), "\n")
}

func (rc *mdResultCollector) open(heading changelog.Heading) {
	rc.prefix = lineStart(heading.Kind())
}

func lineStart(heading changelog.HeadingKind) string {
	switch heading {
	case changelog.ChangeDescription:
		return "- "
	default:
		return ("###"[:(heading + 1)]) + " "
	}
}

func (rc *mdResultCollector) close(heading changelog.Heading) {
}

func (rc *mdResultCollector) setCollection() {
}

func (rc *mdResultCollector) set(value string) {
	rc.result.WriteString(rc.prefix)
	rc.result.WriteString(value)
	rc.result.WriteString("\n")
	rc.prefix = ""
}

func (rc *mdResultCollector) setField(name string, value string) {
	if name == "title" {
		rc.result.WriteString(rc.prefix)
		rc.result.WriteString(value)
		rc.result.WriteString("\n")
		rc.prefix = ""
	}
}

func (rc *mdResultCollector) array(name string) {
}
