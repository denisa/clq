package output

import (
	"fmt"

	"github.com/denisa/clq/internal/changelog"
)

// OutputFormat exposes to the rest of the application the plugin mechanism
// through which multiple output formats are supported.
type OutputFormat interface {
	// Result a string representation of the query result.
	Result() string
	// Open is called when a heading is first met.
	Open(heading changelog.Heading)
	// Close is called when all of a heading’s children have been visited.
	Close(heading changelog.Heading)
	// SetCollection lets the OutputFormat knows that the query will produce a
	// collection of results.
	SetCollection()
	Set(value string)
	SetField(name string, value string)
	Array(name string)
}

func NewOutputFormat(formatName string) (OutputFormat, error) {
	switch formatName {
	case "json":
		return &jsonResultCollector{}, nil
	case "md":
		return &mdResultCollector{}, nil
	default:
		return nil, fmt.Errorf("Unrecognized output format %q. Supported format: \"json\", \"md\"", formatName)
	}
}
