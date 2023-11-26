package output

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

func (rc *mdResultCollector) Open(heading changelog.Heading) {
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

func (rc *mdResultCollector) Close(heading changelog.Heading) {
}

func (rc *mdResultCollector) SetCollection() {
}

func (rc *mdResultCollector) Set(value string) {
	rc.result.WriteString(rc.prefix)
	rc.result.WriteString(value)
	rc.result.WriteString("\n")
	rc.prefix = ""
}

func (rc *mdResultCollector) SetField(name string, value string) {
	if name == "title" {
		rc.result.WriteString(rc.prefix)
		rc.result.WriteString(value)
		rc.result.WriteString("\n")
		rc.prefix = ""
	}
}

func (rc *mdResultCollector) Array(name string) {
}
