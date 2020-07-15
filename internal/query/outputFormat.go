package query

import (
	"fmt"

	"github.com/denisa/clq/internal/changelog"
)

// OutputFormat exposes to the rest of the application the plugin mechanism
// through which multiple output formats are supported.
type OutputFormat interface {
	// Result a string representation of the query result.
	Result() string
	// open is called when a heading is first met.
	open(heading changelog.Heading)
	// Close is called when all of a heading’s children have been visited.
	close(heading changelog.Heading)
	// setCollection lets the OutputFormat kows that the query will produce a
	// collection of results.
	setCollection()
}

// the resultCollector collects the projected fields to be formatted by the
// OutputFormat.
type resultCollector interface {
	set(value string)
	setField(name string, value string)
	array(name string)
}

// a project function projects the desired part of the heading in the resultCollector.
type project func(rc resultCollector, heading changelog.Heading)

// projections is a base type for all queries
type projections struct {
	enter, exit project
	collection  bool
}

func newOutputFormat(formatName string) (OutputFormat, error) {
	switch formatName {
	case "json":
		return &jsonResultCollector{}, nil
	case "md":
		return &mdResultCollector{}, nil
	default:
		return nil, fmt.Errorf("Unrecognized output format %q. Supported format: \"json\", \"md\"", formatName)
	}
}
